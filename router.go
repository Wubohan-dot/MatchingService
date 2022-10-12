package main

import (
	"MatchingService/biz"
	"net/http"
)

// response query:  “==” equal, “!=” not equal, “$=” equal (case insensitive), “&=” contain (the query term is a substring of the data cell)
// also, some signal should be transferred, including :& -> %26 and $ -> %24
// some cases:
// 127.0.0.1:9527/?query=C == "c1" or C %26= "c"
// 127.0.0.1:9527/?query=C == "c1" or C %26=
// 127.0.0.1:9527/?query=C"c1" or C %26= "c"
// 127.0.0.1:9527/?query=C == "c1" or C %26= c
// 127.0.0.1:9527/?query=C == "c1" or A %26= "c"
// 127.0.0.1:9527/?query=C == "c1" or A %26= "a"
func main() {
	http.HandleFunc("/", biz.IndexHandler)
	http.ListenAndServe(":9527", nil)
}
