package service

import (
	"encoding/csv"
	"errors"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func getUnion(choices map[int]struct{}, choices2 map[int]struct{}) map[int]struct{} {
	possibleChoices := make(map[int]struct{})
	for k, _ := range choices {
		possibleChoices[k] = struct{}{}
	}
	for k, _ := range choices2 {
		possibleChoices[k] = struct{}{}
	}
	return possibleChoices
}

func getIntersection(choices map[int]struct{}, choices2 map[int]struct{}) map[int]struct{} {
	possibleChoices := make(map[int]struct{})
	for k, _ := range choices {
		if _, ok := choices2[k]; ok {
			possibleChoices[k] = struct{}{}
		}
	}
	return possibleChoices
}

// separateQueries:
// queries like: C1 == "A" or C2 %26= "B"
// to
// [[C1,==,A,and][C2,&=,B,or]]
func separateQueries(queries string) ([][]string, error) {
	words, err := checkIfQueriesValid(queries)
	if err != nil {
		log.Printf("[seperateQueries] queries not invalid: %v", queries)
		return nil, errors.New("invalid query")
	}
	queryArr := make([][]string, 0)
	for i := 0; i < len(words); i += 4 {
		queryArr = append(queryArr, words[i:i+4])
	}
	return queryArr, nil
}

// C1 == "A" or C2 %26= "B"
func checkIfQueriesValid(queries string) ([]string, error) {
	words := strings.Split(queries, " ")
	// if words is empty, it needs to be handled specially to avoid expose non-exist "and"
	if len(words) == 0 {
		errMsg := "empty query"
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}
	// add an extra "and" at head to make it 4-circle
	words = append([]string{"and"}, words...)
	log.Printf("%v, words: %v", len(words), words)
	for i, word := range words {
		switch i % 4 {
		case 0:
			if word != "and" && word != "or" {
				errMsg := "wrong query near " + word
				log.Printf(errMsg)
				return nil, errors.New(errMsg)
			}
		case 1:
			continue
		case 2:
			if _, ok := operator2Selector[word]; !ok {
				errMsg := "wrong query near " + word
				log.Printf(errMsg)
				return nil, errors.New(errMsg)
			}
		case 3:
			tmp, err := checkIfWordsWithQuotation(word)
			if err != nil {
				return nil, errors.New(err.Error())
			}
			words[i] = tmp
		}
	}
	// if the words cannot be divided by 4, then something near tail wrong
	if len(words)%4 != 0 {
		errMsg := "wrong query near " + words[len(words)-1]
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}
	return words, nil
}

func checkIfWordsWithQuotation(word string) (string, error) {
	n := len(word)
	if n < 2 {
		errMsg := "wrong query near " + word
		log.Printf(errMsg)
		return "", errors.New(errMsg)
	}
	if word[0] != '"' || word[n-1] != '"' {
		errMsg := "wrong query near " + word
		errMsg += " ,need quotation"
		return "", errors.New(errMsg)
	}
	return word[1 : n-1], nil
}

func GenerateFileNameWithRandomIdentifier() string {
	filename := "./output-" + strconv.FormatInt(rand.Int63(), 10)
	filename += ".csv"
	return filename
}

func CSVFormatResp(resp [][]string) []byte {
	var s string
	// [239 187 191 65 44 66 44 67 10 97 49 44 98 49 44 99 49 10]
	for _, row := range resp {
		for _, data := range row[:len(row)-1] {
			s += data
			s += ","
		}
		s += row[len(row)-1]
		s += string(rune(10))
	}
	return []byte(s)
}

func MakeTmpFile(resp [][]string) (string, error) {
	filename := GenerateFileNameWithRandomIdentifier()

	xlsFile, fErr := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0766)
	if fErr != nil {
		log.Println("Export:created excel file failed ==", fErr)
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
