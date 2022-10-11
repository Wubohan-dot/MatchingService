package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
		return
	}

	filename, err := makeTmpFile(resp)

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

func makeTmpFile(resp [][]string) (string, error) {
	filename := "./output-" + strconv.FormatInt(rand.Int63(), 10)
	filename += ".csv"

	xlsFile, fErr := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0766)
	if fErr != nil {
		fmt.Println("Export:created excel file failed ==", fErr)
		return "", errors.New("fail to generate file")
	}
	defer xlsFile.Close()

	xlsFile.WriteString("\xEF\xBB\xBF")

	wStr := csv.NewWriter(xlsFile)
	wStr.Write(resp[0])

	for _, s := range resp[1:] {
		wStr.Write(s)
	}
	wStr.Flush()
	return filename, nil
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":9527", nil)
}
