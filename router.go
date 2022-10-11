package main

import (
	"fmt"
	"net/http"
)

// 127.0.0.1:9527/?query=C1 == "A" or C2 %26= "B"
func indexHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte(string("post")))
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println(r.Form["query"])
	if len(r.Form["query"]) == 0 {
		w.Write([]byte("need a query"))
	} else {
		w.Write([]byte(r.Form["query"][0]))
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":9527", nil)
}
