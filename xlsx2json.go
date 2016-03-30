//
// xlsx2json package wraps the github.com/tealag/xlsx package (used under a BSD License) and  a fork of Robert Krimen's Otto
// Javascript engine (under an MIT License) providing an scriptable xlsx2json exporter, explorer and importer utility.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
//
// Copyright (c) 2016, R. S. Doiel
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package xlsx2json

import (
	"encoding/json"
	"fmt"

	// 3rd party packages
	"github.com/robertkrimen/otto"
	"github.com/tealeg/xlsx"
)

// Version is the library and utilty version number
const Version = "0.0.1"

type jsResponse struct {
	Path   string                 `json:"path"`
	Source map[string]interface{} `json:"source"`
	Error  string                 `json:"error"`
}

// Run runs the xlsx2json transform with optional JavaScript support.
// Continued processing can be achieved with subsequent calls to
// the JS VM. It returns the VM, an array of JSON encoded blobs and error.
func Run(inputFilename string, sheetNo int, jsFilename, jsCallback string) (*otto.Otto, []string, error) {
	var (
		xlFile *xlsx.File
		vm     *otto.Otto
		err    error
		output []string
	)

	jsMap := false
	vm, err = NewJavaScriptVM([]string{jsFilename})
	if err != nil {
		return nil, output, err
	}
	if jsFilename != "" && jsCallback != "" {
		jsMap = true
	}

	// Read from the given file path
	xlFile, err = xlsx.OpenFile(inputFilename)
	if err != nil {
		return vm, output, fmt.Errorf("Can't open %s, %s", inputFilename, err)
	}

	for i, sheet := range xlFile.Sheets {
		if sheetNo == i {
			columnNames := []string{}

			for rowNo, row := range sheet.Rows {
				jsonBlob := map[string]string{}
				for colNo, cell := range row.Cells {
					if rowNo == 0 {
						s, _ := cell.String()
						columnNames = append(columnNames, s)
					} else {
						// Build a map and render it out
						if colNo < len(columnNames) {
							s, _ := cell.String()
							jsonBlob[columnNames[colNo]] = s
						} else {
							k := fmt.Sprintf("column_%d", colNo+1)
							columnNames = append(columnNames, k)
							s, _ := cell.String()
							jsonBlob[k] = s
						}
					}
				}
				if rowNo > 0 {
					src, err := json.Marshal(jsonBlob)
					if err != nil {
						return vm, output, fmt.Errorf("Can't render JSON blob, %s", err)
					}
					if jsMap == true {
						// We're eval the callback from inside a closure to be safer
						js := fmt.Sprintf("(function(){ return %s(%s);}())", jsCallback, src)
						jsValue, err := vm.Eval(js)
						if err != nil {
							return vm, output, fmt.Errorf("row: %d, Can't run %s, %s", rowNo, jsFilename, err)
						}
						val, err := jsValue.Export()
						if err != nil {
							return vm, output, fmt.Errorf("row: %d, Can't convert JavaScript value %s(%s), %s", rowNo, jsCallback, src, err)
						}
						src, err = json.Marshal(val)
						if err != nil {
							return vm, output, fmt.Errorf("row: %d, src: %s\njs returned %v\nerror: %s", rowNo, js, jsValue, err)
						}
						response := new(jsResponse)
						err = json.Unmarshal(src, &response)
						if err != nil {
							return vm, output, fmt.Errorf("row: %d, do not understand response %s, %s", rowNo, src, err)
						}
						if response.Error != "" {
							return vm, output, fmt.Errorf("row: %d, %s", rowNo, response.Error)
						}
						// Now re-package response.Source into a JSON blob
						src, err = json.Marshal(response.Source)
						if err != nil {
							return vm, output, fmt.Errorf("row: %d, %s", rowNo, err)
						}
					}
					output = append(output, string(src))
				}
			}
		}
	}
	return vm, output, nil
}
