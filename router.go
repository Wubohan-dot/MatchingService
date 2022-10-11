package main

import (
	"log"
	"net/http"
)

func checkIfValid(r *http.Request) error {
	return nil
}

// 127.0.0.1:9527/?query=C1 == "A" or C2 %26= "B"
func indexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	err := checkIfValid(r)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
	matcherInit()
	resp, err := matcher.matchWithQueries(r.Form["query"][0])
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
