/*
 *  A simplistic example of processing an Excel row with JavaScript
 */
function slugify(s) {
    return 'testdata/' + encodeURI(s.replace(" ", "_")) + '.json';
}

function row2object(row) {
    if (row.Name == undefined) {
        return {"Path":"", "Source": "", "Error": "Missing Name property"}
    }
    return {
        "Path": slugify(row.Name),
        "Source": row,
        "Error": null
    };
}

