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
