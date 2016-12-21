/* counter.js - Add column Counter saving i in the JSON output. */
var i = 0;

/* 
 * callback is the default name looked for when processing.
 *  the command line option -callback lets you used a different name.
 */
function callback(row) {
    i += 1;
    if (i > 10) {
        /* Stop if processing more than 10 rows. */
        return {"error": "too many rows..."}
    }
    /* Add a counter column and save the current value of i */
    row.Counter = i
    return {
        "path": "data/" + i + ".json",
        "source": row,
        "error": ""
    }
}
