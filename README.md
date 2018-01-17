
    this repository has been depreciated and merged with https://github.com/caltechlibrary/datatools

# xlsx2json

This is a simple command line utility for converting XML based Excel files to useful JSON objects.  It uses Geoff Teale's [xlsx](https://github.com/tealeg/xlsx) golang package alone with Robert Krimen's [Otto JavaScript interpretor](https://github.com/robertkrimen/otto).

The basic idea is that xlsx2json reads an Excel file and outputs a JSON expression for each row of the spreadsheet. The headings of the columns become the key names and the value is cell's content.  If you provide a mapping functions written in JavaScript that will be use to process a row into a JSON object.

Like most Unix command line utilities xlsx2json will write to standard output (or standard error if there is a problem). Input is always assumed to be an Excel file's content, output is JSON blobs (one row is one blob).

## Basic usage

```
    xlsx2json myfile.xlsx
```

This will write a series of JSON blobs to standard out (assuming no error, errors are written to standard error).

Because excel spreadsheets are typically complete documents the entire sheet(s) is read in before being converted and written out.

## Controlling output

You can control the out through providing a JavaScript mapping function.  The utility will execute a function a specified callback at the each row encountered. The entire row is always processed in one step.  This allows you to make calculations and before the JSON output is rendered.  The JavaScript callback should accept a row object as a parameter and return an object structured with a Path attribute, Source attribute (a JSON expression of the rendered content) and an Error attribute. If Path is empty then standard output is assumed.  Otherwise _xlsx2json_ will assume you want to write the JSON content in Source to the filepath provided. If Error is not an empty string then that will be written to standard Error and no output will be written to Path (or standard out) for that row.

Here's a basic example for writing _myfile.xlsx_ as a series of JSON files mapped with _myobjects.js_

```
    xlsx2json -js row2obj.js -callback row2obj myfile.xlsx
```

_row2obj.js_ might look something like this--

```JavaScript
    /*
     *  A simplistic example of processing an Excel row with JavaScript
     *  The excel file has a Name, Email, Age and Journeyman Status columns.
     *  This script will generate an object per row with the following fields--
     *
     *  + Name (string)
     *  + Email (string)
     *  + Age (numeric value)
     *  + Journeyman (boolean)
     *
     */
    
    var i = 0;
    
    // Use the "Name" column to determine the filename by making it file system name friendly
    function slugify(s) {
        return 'testout/' + encodeURI(s.replace(" ", "_")) + '.json';
    }
    
    // Given a row object write the contents to a JSON blob file.
    // @param row - row comes in as an object with column names for property values
    // @return an object with a path, source and error values. path and error are strings
    // source is the modified object representing the updated object structure used for the final
    // JSON output.
    function row2obj(row) {
        if (row.Name === undefined) {
            return {"Path":"", "Source": "", "Error": "Missing Name property"}
        }
        // Add the count
        i++;
        row.Count = i;
    
        // Convert "%d" into a numeric value
        a = parseInt(row.Age, 10);
        row.Age = a;
        
        // Convert "1" and "0" to true and false
        if (row["Journeyman Status"] !== undefined && row["Journeyman Status"] == "1") {
            row.Journeyman = true;
        } else {
            row.Journeyman = false;
        }
        delete row["Journeyman Status"];
        return {
            "path": slugify(row.Name),
            "source": row,
            "error": ""
        };
    }
```

## Installation

The golang compiler (version 1.7.3 or better) needs to be installed.

Clone the [source repository](https://github.com/rsdoiel/xlsx2json) on Github.
Then change to the directory where your cloned repository is.

```
    git clone https://github.com/rsdoiel/xlsx2json
    cd xlsx2json
```

Next you need to have the xlsx and otto packages available. They are "go get"-able.

```
    go get -u github.com/tealeg/xlsx
    go get -u github.com/robertkrimen/otto
```

Then you can run

```
    go build -o bin/xlsx2json cmds/xlsx2json/xlsx2json.go
```

To create the xlsx2json command utility.

You can also just run "make" as their is a simple makefile for building and
installation.

## Installation

_xlsx2json_ can be installed with the *go get* command.

```
    go get github.com/rsdoiel/xlsx2json/...
```

