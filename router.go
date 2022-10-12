package main

import (
	"MatchingService/biz"
	"net/http"
)

func main() {
	http.HandleFunc("/", biz.IndexHandler)
	http.ListenAndServe(":9527", nil)
}
