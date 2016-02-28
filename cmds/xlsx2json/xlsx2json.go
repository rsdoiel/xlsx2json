//
// xlsx2json package wraps the github.com/tealag/xlsx package (used under a BSD License) and  a fork of Robert Krimen's Otto
// Javascript engine (under an MIT License) providing an scriptable xlsx2json exporter, explorer and importer utility.
//
//
// Overview: A command line utility designed to take a XML based Excel file
// and turn each row into a JSON blob. The JSON blob returned for
// each row can be processed via a JavaScript callback allowing for
// flexible translations for spreadsheet to JSON output.
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
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// 3rd Party packages
	"github.com/chzyer/readline"

	// My packages
	"github.com/rsdoiel/xlsx2json"
)

var (
	showhelp      bool
	sheetNo       int
	inputFilename *string
	jsFilename    *string
	jsCallback    *string
)

func usage() {
	fmt.Println(`

 USAGE: xlsx2json [OPTIONS] EXCEL_FILENAME

 OVERVIEW

 Read a .xlsx file and return each row as a JSON object (or array of objects).
 If a JavaScript file and callback name are provided then that will be used to
 generate the resulting JSON object per row.

 JAVASCRIPT

 The callback function in JavaScript should return an object that looks like

     {"path": ..., "source": ..., "error": ...}

 The "path" property should contain the desired filename to use for storing
 the JSON blob. If it is empty the output will only be displayed to standard out.

 The "source" property should be the final version of the object you want to
 turn into a JSON blob.

 The "error" property is a string and if the string is not empty it will be
 used as an error message and cause the processing to stop.

 A simple JavaScript Examples:

    // Counter i is used to name the JSON output files.
    var i = 0;

    // callback is the default name looked for when processing.
    // the command line option -callback lets you used a different name.
    function callback(row) {
        i += 1;
        if (i > 10) {
            // Stop if processing more than 10 rows.
            return {"error": "too many rows..."}
        }
        return {
            "path": "data/" + i + ".json",
            "source": row,
            "error": ""
        }
    }


 OPTIONS
`)
	flag.PrintDefaults()
	fmt.Println(`

 Examples

    xlsx2json myfile.xlsx

    xlsx2json -js row2obj.js -callback row2obj myfile.xlsx

	xlsx2json -repl myfile.xlsx

`)
	os.Exit(0)
}

func init() {
	flag.BoolVar(&showhelp, "h", false, "display this help message")
	flag.BoolVar(&showhelp, "help", false, "display this help message")
	flag.BoolVar(&showVersion, "v", false, "display version information")
	flag.BoolVar(&xlsx2json.UseRepl, "i", false, "Run with an interactive repl")
	flag.IntVar(&sheetNo, "sheet", 0, "Specify the number of the sheet to process")
	jsFilename = flag.String("js", "", "The name of the JavaScript file containing callback function")
	jsCallback = flag.String("callback", "callback", "The name of the JavaScript function to use as a callback")
}

func main() {
	var (
		inputFilename string
	)
	flag.Parse()

	if showhelp == true {
		usage()
	}
	if showVersion == true {
		fmt.Printf("Version %s\n", xlsx2json.Version)
	}

	args := flag.Args()
	if len(args) > 0 {
		inputFilename = args[0]
	}
	if inputFilename == "" {
		// Read Excel file from standard
		log.Fatalf("Need to provide an xlsx file for input, -i")
	}
	vm, output, err := xlsx2json.Run(inputFilename, sheetNo, *jsFilename, *jsCallback)
	if err != nil {
		log.Fatal(err)
	}
	if xlsx2json.UseRepl == true {
		//FIXME: merge output as an array
		vm.Object(fmt.Sprintf(`Spreadsheet = [%s]`, strings.Join(output, ",")))
		rl, err := readline.New("> ")
		if err != nil {
			log.Fatal(err)
		}
		defer rl.Close()
		for xlsx2json.UseRepl == true {
			jsSrc, err := rl.Readline()
			if err != nil { // io.EOF, readline.ErrInterrupt
				break
			}
			if len(strings.Trim(jsSrc, " ")) > 0 {
				if script, err := vm.Compile("repl", jsSrc); err != nil {
					fmt.Printf("Compile error, %s\n", err)
				} else {
					out, err := vm.Eval(script)
					switch {
					case err != nil:
						fmt.Printf("Runtime error, %s\n", err)
					default:
						if xlsx2json.UseRepl == true {
							fmt.Println(out.String())
						}
					}
				}
			}
		}
		//FIXME: update output to reflect the contents of Spreadsheet from the JS VM
		if value, err := vm.Get("Spreadsheet"); err == nil {
			data, err := value.Export()
			if err != nil {
				log.Fatal(err)
			}
			src, err := json.Marshal(data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s", src)
			os.Exit(0)
		} else {
			log.Fatal(err)
		}
	}
	fmt.Printf(`[%s]`, strings.Join(output, ","))
}
