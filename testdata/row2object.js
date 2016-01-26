/*
 *  A simplistic example of processing an Excel row with JavaScript
 */

// Use the "Name" column to determine the filename by making it file system name friendly
function slugify(s) {
    return 'testdata/' + encodeURI(s.replace(" ", "_")) + '.json';
}

// Given a row object write the contents to a JSON blob file.
function row2object(row) {
    if (row.Name == undefined) {
        return {"Path":"", "Source": "", "Error": "Missing Name property"}
    }
    return {
        "Path": slugify(row.Name),
        "Source": row,
        "Error": ""
    };
}

