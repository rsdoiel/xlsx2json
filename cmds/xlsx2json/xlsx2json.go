//
// Package aspace is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2016
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/robertkrimen/otto"
	"github.com/tealeg/xlsx"
)

var (
	help          bool
	asArray       bool
	inputFilename *string
	jsFilename    *string
	jsCallback    *string
)

type jsResponse struct {
	Path   string
	Source map[string]interface{}
	Error  string
}

func usage() {
	fmt.Println(`
 USAGE: xlsx2json [OPTIONS] EXCEL_FILENAME

 OVERVIEW

 Read a .xlsx file and return each row as a JSON object (or array of objects)

 OPTIONS
`)
	flag.PrintDefaults()
	fmt.Println(`

 Examples

    xlsx2json -i myfile.xlsx -as-array

    xlsx2json -i myfile.xlsx -js obj2row.js -callback obj2row

`)
	os.Exit(0)
}

func init() {
	flag.BoolVar(&help, "h", false, "display this help message")
	flag.BoolVar(&help, "help", false, "display this help message")
	flag.BoolVar(&asArray, "as-array", false, "Write the JSON blobs output as an array")
	inputFilename = flag.String("i", "", "Read the Excel file from this name")
	jsFilename = flag.String("js", "", "The name of the JavaScript file containing callback function")
	jsCallback = flag.String("callback", "", "The name of the JavaScript function to use as a callback")
}

func main() {
	var (
		xlFile   *xlsx.File
		vm       *otto.Otto
		jsSource []byte
        jsScript *otto.Script
	)
	flag.Parse()

	if help == true {
		usage()
	}

    args := flag.Args()
    if len(args) > 0 {
        *inputFilename = args[0]
    }
	if *inputFilename == "" {
		// Read Excel file from standard
		log.Fatalf("Need to provide an xlsx file for input, -i")
	}
	// Read from the given file path
	xlFile, err := xlsx.OpenFile(*inputFilename)
	if err != nil {
		log.Fatalf("Can't open %s, %s", *inputFilename, err)
	}
	jsMap := false
	if *jsFilename != "" {
		fname := fmt.Sprintf("%s", *jsFilename)
		jsSource, err = ioutil.ReadFile(fname)
		if err != nil {
			log.Fatalf("Can't read JavaScript file %s, %s", fname, err)
		}
		vm = otto.New()
		_, err := vm.Run(jsSource)
		if err != nil {
			log.Fatalf("Can't run %s, %s", *jsFilename, err)
		}
		jsMap = true
        jsScript, err = vm.Compile(*jsFilename, jsSource)
        if err != nil {
            log.Fatalf("Can't compile %s, %s", *jsFilename, err)
        }
		//FIXME: add JS wrapped Golang packages
        // Define any functions, will evaluate each row with vm.Eval()
        vm.Run(jsScript)
	}

	for _, sheet := range xlFile.Sheets {
		columnNames := []string{}
		if asArray == true {
			fmt.Println("[")
		}
		for rowNo, row := range sheet.Rows {
			if asArray == true && rowNo > 1 {
				fmt.Printf(", ")
			}
			jsonBlob := map[string]string{}
			for colNo, cell := range row.Cells {
				if rowNo == 0 {
					columnNames = append(columnNames, cell.String())
				} else {
					// Build a map and render it out
					if colNo < len(columnNames) {
						jsonBlob[columnNames[colNo]] = cell.String()
					} else {
						k := fmt.Sprintf("column_%d", colNo+1)
						columnNames = append(columnNames, k)
						jsonBlob[k] = cell.String()
					}
				}
			}
			if rowNo > 0 {
				src, err := json.Marshal(jsonBlob)
				if err != nil {
					log.Fatalf("Can't render JSON blob, %s", err)
				}
				if jsMap == true {
					// We're eval the callback from inside a closure to be safer
					js := fmt.Sprintf("(function(){ return %s(%s);}())", *jsCallback, src)
                    jsValue, err := vm.Eval(js)
					if err != nil {
						log.Fatalf("row: %d, Can't run %s, %s", rowNo, *jsFilename, err)
					}
					val, err := jsValue.Export()
					if err != nil {
						log.Fatalf("row: %d, Can't convert JavaScript value %s(%s), %s", rowNo, *jsCallback, src, err)
					}
					src, err = json.Marshal(val)
					if err != nil {
						log.Fatalf("row: %d, src: %s\njs returned %v\nerror: %s", rowNo, js, jsValue, err)
					}
					response := new(jsResponse)
					err = json.Unmarshal(src, &response)
					if err != nil {
						log.Fatalf("row: %d, do not understand response %s, %s", rowNo, src, err)
					}
					if response.Error != "" {
						log.Fatalf("row: %d, %s", rowNo, response.Error)
					}
                    // Now re-package response.Source into a JSON blob
                    src, err = json.Marshal(response.Source)
					if err != nil {
                        log.Fatalf("row: %d, %s", rowNo, err)
					}
					if response.Path != "" {
						d := path.Dir(response.Path)
						if d != "." {
							os.MkdirAll(d, 0775)
						}
						ioutil.WriteFile(response.Path, src, 0664)
					}
				}
				fmt.Printf("%s\n", src)
			}
		}
		if asArray == true {
			fmt.Println("]")
		}

	}
}
