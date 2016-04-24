
[![Go Report Card](http://goreportcard.com/badge/rsdoiel/xlsx2json)](http://goreportcard.com/report/rsdoiel/xlsx2json)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)


# xlsx2json

This is a simple command line utility for converting Excel Workbook files (in the .xlsx XML based Excel file format) to useful JSON objects.  It uses Geoff Teale's [xlsx](https://github.com/tealeg/xlsx) golang package alone with Robert Krimen's [Otto JavaScript interpretor](https://github.com/robertkrimen/otto).

The basic idea is that xlsx2json reads an Excel file and outputs a JSON expression with each workbook an object with the property name corresponds to the sheet name and the property value us a 2d array of string. 

```
    {
        "Sheet1": [
            ["Heading 1", "Heading 2"],
            ["one", "two"]
        ]
    }
```

You can also apply a callback process the workbook before displaying the output (this allowing for alternative output organizations).

Like most Unix command line utilities xlsx2json will write to standard output (or standard error if there is a problem). Input is always assumed to be an Excel file's content, output is JSON object representing the entire Workbook.

## Basic usage

```
    xlsx2json myfile.xlsx
```

This will write a series of JSON blobs representing the workbook to standard out (assuming no error, errors are written to standard error).

The same writing the JSON objects as an array of objects to myfile.json.

```
    xlsx2json myfile.xlsx > myfile.json
```

Converting the JSON output to something else with a JavaScript callback. In this case the callback is named newFormat and it is defined
in the JavaScript file named myformats.js

```
    xlsx2json -callback newFormat myfile.xlsx myformats.js
```

Because excel spreadsheets are typically complete documents the entire sheet(s) is read in before being converted and written out.


## Installation

The golang compiler (version 1.6 or better) needs to be installed.

```
    go get github.com/rsdoiel/xlsx2json/...
```

Or grab the sourcecode by cloning [github.com/rsdoiel/xlsx](https://github.com/rsdoiel/xlsx2json).

```
    git clone https://github.com/rsdoiel/xlsx2json
    cd xlsx2json
```

If you "go get" the ostdlib package you'll get all the other dependencies such as xlsx, readline, color and otto packages. 

```
    go get github.com/caltechlibrary/ostdlib
```

## Installation

_xlsx2json_ can be installed with the *go get* command.

```
    go get github.com/rsdoiel/xlsx2json/...
```

