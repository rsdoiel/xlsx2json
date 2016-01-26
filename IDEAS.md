
The JS engine needs to support a few Golang standard library calls

+ ioutil
+ os
+ http client code for Get, Post, Put and Delete
+ path 
+ url

The JS engine to create the VM, create these default objects, compile the supplied script
and then for each row of the spreadsheet eval the callback function and process the results.


