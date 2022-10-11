package main

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
)

type Matcher struct {
	dict       map[string][]string
	columnNum  int
	choicesNum int
}

var matcher Matcher

func init() {
	matcher.dict = make(map[string][]string)
	matcher.getDict("E:\\MatchingService\\dict.csv")
}

func (m *Matcher) getDict(fileName string) error {
	records, err := m.readFile(fileName)
	if err != nil {
		log.Printf("[getDict] err:%v", err)
		return err
	}
	err = m.formatToDict(records)
	if err != nil {
		log.Printf("[getDict] err:%v", err)
		return err
	}
	return nil
}

func (m *Matcher) readFile(fileName string) ([][]string, error) {
	opencast, err := os.Open(fileName)
	if err != nil {
		log.Printf("[readDict] fail to open file: %v, error: %v", fileName, err)
		return nil, err
	}
	defer opencast.Close()

	reader := csv.NewReader(opencast)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("[readDict] error: %v", err)
		return nil, err
	}

	log.Printf("[readDict] success, record: %+v with %v rows and %v columns", records, len(records), len(records[0]))
	return records, nil
}

func (m *Matcher) formatToDict(records [][]string) error {
	err := m.checkIfRecordsValid(records)
	if err != nil {
		return errors.New("illegal Records")
	}

	m.columnNum = len(records[0])
	m.choicesNum = len(records) - 1

	for i := 0; i < m.columnNum; i++ {
		title := records[0][i]
		for j := 0; j < m.choicesNum; j++ {
			m.dict[title] = append(m.dict[title], records[j+1][i])
		}
	}
	log.Println(m.dict)
	log.Printf("[formatToDict] success with m.columnNum:%v and choicesNum:%v", m.columnNum, m.choicesNum)
	return nil
}

func (m *Matcher) checkIfRecordsValid(records [][]string) error {
	// check if there are at least two rows: one for title and one for content
	if len(records) < 2 {
		log.Printf("[checkIfRecordsValid] lack of content in .csv file")
		return errors.New("lack of content in .csv file")
	}
	// check if each row has the same length
	for i, row := range records[:len(records)-1] {
		if len(row) != len(records[i+1]) {
			log.Printf("[checkIfRecordsValid] have %v elements in row %v while %v elements in row %v", len(row), i, len(records[i+1]), i+1)
			return errors.New("wrong contents in records")
		}
	}
	// check the title line only contains A-Z, a-z, 0-9(case sensitive)
	for _, title := range records[0] {
		for _, char := range title {
			if (char < '0' || char > '9') && (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
				log.Printf("[checkIfRecordsValid] illegal tile")
				return errors.New("illegal title in records")
			}
		}
	}
	return nil
}
