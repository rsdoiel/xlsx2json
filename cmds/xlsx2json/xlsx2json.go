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
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	// 3rd Party packages
	"github.com/robertkrimen/otto"

	// Caltech Library packages
	"github.com/caltechlibrary/ostdlib"

	// My packages
	"github.com/rsdoiel/cli"
	"github.com/rsdoiel/xlsx2json"
)

var (
	usage = `USAGE: %s [OPTIONS] EXCEL_FILENAME`

	description = ` 
SYNOPSIS

%s reads an workbook .xlsx file and returns each row as a JSON object (or array of objects).
If a JavaScript file and callback name are provided then that will be used to
generate the resulting JSON content.

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

    /* Counter i is used to name the JSON output files. */
    var i = 0;

    /* callback is the default name looked for when processing.
       the command line option -callback lets you used a different name. */
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
`

	examples = `
EXAMPLES

    %s myfile.xlsx

    %s -callback row2obj row2obj.js myfile.xlsx

	%s -i myfile.xlsx
`

	// Standard Options
	showhelp    bool
	showLicense bool
	showVersion bool

	// Application Options
	sheetNo       int
	jsCallback    string
	jsFilename    string
	jsInteractive bool
	xlsxFilename  string
)

func init() {
	// Standard Options
	flag.BoolVar(&showhelp, "h", false, "display help")
	flag.BoolVar(&showhelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showVersion, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")

	// Application Options
	flag.BoolVar(&jsInteractive, "i", false, "Run with an interactive repl")
	flag.IntVar(&sheetNo, "sheet", 0, "Specify the number of the sheet to process")
	flag.StringVar(&jsFilename, "j", "", "JavaScript filename")
	flag.StringVar(&jsFilename, "js", "", "JavaScript filename")
	flag.StringVar(&jsCallback, "c", "", "The name of the JavaScript function to use as a callback")
	flag.StringVar(&jsCallback, "callback", "", "The name of the JavaScript function to use as a callback")
}

func main() {
	var (
		output []string
		err    error
	)
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	cfg := cli.New(appName, "XLSX2JSON",
		fmt.Sprintf(xlsx2json.LicenseText, appName, xlsx2json.Version),
		xlsx2json.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = description
	cfg.OptionsText = "OPTIONS\n"
	cfg.ExampleText = fmt.Sprintf(examples, appName, appName, appName)

	// handle Standard Options
	if showhelp == true {
		fmt.Println(cfg.Usage())
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	for _, fname := range args {
		if strings.HasSuffix(fname, ".js") {
			jsFilename = fname
		}
		if strings.HasSuffix(fname, ".xlsx") {
			xlsxFilename = fname
		}
	}

	if len(xlsxFilename) == 0 && jsInteractive == false {
		fmt.Fprintf(os.Stderr, "Missing xlsx filename, see %s -help\n", appName)
		os.Exit(1)
	}

	vm := otto.New()
	js := ostdlib.New(vm)
	js.AddExtensions()

	if jsFilename != "" {
		// Check to see if we need to use the default callback
		if len(jsCallback) == 0 {
			jsCallback = "callback"
		}
		if err := js.Run(jsFilename); err != nil {
			log.Fatalf("%s", err)
		}
	}

	if len(xlsxFilename) > 0 {
		output, err = xlsx2json.Run(js, xlsxFilename, sheetNo, jsCallback)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	// Join the preformatted strings into a JSON array
	src := fmt.Sprintf("[%s]", strings.Join(output, ","))
	if jsInteractive == true {
		js.AddHelp()
		js.AddAutoComplete()
		js.PrintDefaultWelcome()
		js.Eval(fmt.Sprintf("Spreadsheet = %s;", src))
		js.Repl()
	} else {
		fmt.Print(src)
	}
}
