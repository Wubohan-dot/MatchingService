package main

import (
	"encoding/csv"
	"log"
	"os"
)

type Matcher struct {
	dict map[string][]string
}

func (m *Matcher) readDict(fileName string) error {
	opencast, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("[readDict] fail to open file: %v, error: %v", fileName, err)
		return err
	}
	defer opencast.Close()

	reader := csv.NewReader(opencast)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("[readDict] error: %v", err)
		return err
	}

	err = checkIfDictValid(record)
	if err != nil {
		log.Fatalf("[readDict] not valid: %v", err)
		return err
	}

	log.Printf("[readDict] success, record: %+v with %v rows and %v columns", record, len(record), len(record[0]))
	return nil
}

func checkIfDictValid([][]string) error {
	return nil
}

var matcher Matcher

func init() {
	matcher.readDict("E:\\MatchingService\\dict.csv")
}
