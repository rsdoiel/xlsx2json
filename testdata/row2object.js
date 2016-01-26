/*
 *  A simplistic example of processing an Excel row with JavaScript
 */

// Use the "Name" column to determine the filename by making it file system name friendly
function slugify(s) {
    return 'testout/' + encodeURI(s.replace(" ", "_")) + '.json';
}

// Given a row object write the contents to a JSON blob file.
function row2object(row) {
    if (row.Name === undefined) {
        return {"Path":"", "Source": "", "Error": "Missing Name property"}
    }
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
        "Path": slugify(row.Name),
        "Source": row,
        "Error": ""
    };
}
