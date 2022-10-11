package main

import (
	"errors"
	"log"
	"net/http"
)

func checkIfValid(r *http.Request) error {
	r.ParseForm()
	if len(r.Form) < 1 {
		return errors.New("query is empty")
	}
	return nil
}

// 127.0.0.1:9527/?query=C == "c1" or C %26= "c"
// &	%26 contains
// $	%24 insensitive case
// 127.0.0.1:9527/?query=C == "c1" or C %26=
// 127.0.0.1:9527/?query=C"c1" or C %26= "c"
// 127.0.0.1:9527/?query=C == "c1" or C %26= c
// 127.0.0.1:9527/?query=C == "c1" or A %26= "c"
// 127.0.0.1:9527/?query=C == "c1" or A %26= "a"
func indexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	err := checkIfValid(r)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := matcher.matchWithQueries(r.Form["query"][0])
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	for _, res := range resp {
		for _, re := range res {
			w.Write([]byte(re))
		}
	}

	//w.Write([]byte(r.Form["query"][0]))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":9527", nil)
}
