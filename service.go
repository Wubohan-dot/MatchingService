package main

import (
	"encoding/csv"
	"log"
	"os"
)

type Matcher struct {
	dict map[string][]string
}

func (m *Matcher) readFile(fileName string) ([][]string, error) {
	opencast, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("[readDict] fail to open file: %v, error: %v", fileName, err)
		return nil, err
	}
	defer opencast.Close()

	reader := csv.NewReader(opencast)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("[readDict] error: %v", err)
		return nil, err
	}

	log.Printf("[readDict] success, record: %+v with %v rows and %v columns", records, len(records), len(records[0]))
	return records, nil
}

func (m *Matcher) getDict(fileName string) error {
	records, err := m.readFile(fileName)
	if err != nil {
		log.Fatalf("[getDict] err:%v", err)
		return err
	}
	err = m.formatToDict(records)
	if err != nil {
		log.Fatalf("[getDict] err:%v", err)
		return err
	}
	return nil
}

func (m *Matcher) formatToDict(records [][]string) error {
	return nil
}

var matcher Matcher

func init() {
	matcher.getDict("E:\\MatchingService\\dict.csv")
}
