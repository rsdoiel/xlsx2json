
# xlsx2json

This is a simple command line utility for converting XML based Excel files to useful JSON objects.  It uses Geoff Teale's [xlsx](https://github.com/tealeg/xlsx) golang package alone with Robert Krimen's [Otto JavaScript interpretor](https://github.com/robertkrimen/otto).

The basic idea is that xlsx2json reads and Excel file and outputs a JSON expression for each row of the spreadsheet. The headings of the columns become the key names and the value is cell's content.  If you provide a mapping functions written in JavaScript that will be use to process a row into a JSON object.

Like most Unix command line utilities xlsx2json will read from standard input and write to standard output (or standard error if there is a problem). Input is always assumed to be an Excel file's content, output is JSON blobs (one row is one blob).

## Basic usage

```
    xlsx2json myfile.xlsx
```

This will write a series of JSON blobs to standard out (assuming no error, errors are written to standard error).

If you want to output the JSON blobs as an array of JSON blobs then there is an option --as-array which will result in valid JSON output.

Because excel spreadsheets are typically complete documents the entire sheet(s) is read in before being converted and written out.

## Controlling output

You can control the out through providing a JavaScript mapping function.  The utility will execute a function a specified callback at the each row encountered. The entire row is always processed in one step.  This allows you to make calculations and before the JSON output is rendered.  The JavaScript callback should accept a row object as a parameter and return an object structured with a Path attribute, Source attribute (a JSON expression of the rendered content) and an Error attribute. If Path is empty then standard output is assumed.  Otherwise _xlsx2json_ will assume you want to write the JSON content in Source to the filepath provided. If Error is not an empty string then that will be written to standard Error and no output will be written to Path (or standard out) for that row.

Here's a basic example for writing _myfile.xlsx_ as a series of JSON files mapped with _myobjects.js_

```
    xlsx2json -callback myobjects.js myfile.xlsx
```

_myobjects.js_ might look something like this--

```JavaScript
    function sluggifyName(name) {
        return urlencode(name.replace(" ",""));
    }

    // Write each row to its own *.json file based on a sluggified name
    function myobjects(row) {
       myobj = {
           Name: row["column_1"],
           Email: row["column_2"]
       };
        return { "Path": "data"+sluggifyName(myobj.Name)+".json", "Source": myobj, "Error":""}
    }
```


