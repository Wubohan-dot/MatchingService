package biz

import (
	"MatchingService/service"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	err := checkIfValid(r)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := service.MatcherInstance().MatchWithQueries(r.Form["query"][0])
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	filename := service.GenerateFileNameWithRandomIdentifier()

	content := service.CSVFormatResp(resp)
	log.Printf("%+v", content)

	fileName := url.QueryEscape(filename) // to avoid Chinese
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Write(content)

}

func IndexHandlerWithFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	err := checkIfValid(r)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := service.MatcherInstance().MatchWithQueries(r.Form["query"][0])
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	filename, err := service.MakeTmpFile(resp)

	// https://segmentfault.com/a/1190000020202158
	file, err := os.Open(filename)
	content, err := ioutil.ReadAll(file)
	file.Close()
	fileName := url.QueryEscape(filename) // to avoid Chinese
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Write(content)

	defer func() {
		err = os.Remove(filename)
		if err != nil {
			fmt.Println("remove  excel file failed", err)
		}
	}()
}
