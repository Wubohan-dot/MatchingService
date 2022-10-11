package main

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strings"
)

type Matcher struct {
	originalDict  [][]string
	dict          map[string][]string
	columnNum     int
	choicesNum    int
	invertedIndex map[string][]string
}

var matcher Matcher
var operator2Selector map[string]Selector

func init() {
	matcher.dict = make(map[string][]string)
	matcher.getDictAndOriginalDict("E:\\MatchingService\\dict.csv")
	matcher.getInvertedIndex()
	selectorInit()
}

func selectorInit() {
	operator2Selector = make(map[string]Selector)
	operator2Selector["=="] = new(EqualSelector)
	operator2Selector["!="] = new(NotEqualSelector)
	operator2Selector["&="] = new(ContainSelector)
	operator2Selector["$="] = new(EqualInsensitiveCaseSelector)
}

func (m *Matcher) getDictAndOriginalDict(fileName string) error {
	records, err := m.readFile(fileName)
	if err != nil {
		log.Printf("[getDictAndOriginalDict] err:%v", err)
		return err
	}
	m.originalDict = records
	err = m.formatToDict(records)
	if err != nil {
		log.Printf("[getDictAndOriginalDict] err:%v", err)
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

func (m *Matcher) matchWithQueries(queries string) ([][]string, error) {
	queryArr, err := separateQueries(queries)
	if err != nil {
		log.Printf("[matchWithQueries] cannot separate queries to array: %v\n", queries)
		return nil, errors.New("invalid queries")
	}
	// queryNum := len(queryArr)
	possibleChoices := m.getAllPossibleChoices()
	//  “==” equal, “!=” not equal, “$=” equal (case insensitive), “&=” contain (the query term is a substring of the data cell)
	for _, query := range queryArr {
		selector, ok := operator2Selector[query[1]]
		if !ok {
			return nil, errors.New("wrong operator")
		}
		tmpPossibleChoices := selector.selectWithQuery(&matcher, query)
		if query[3] == "and" {
			possibleChoices = getIntersection(tmpPossibleChoices, possibleChoices)
		} else {
			possibleChoices = getUnion(tmpPossibleChoices, possibleChoices)
		}

	}
	resp := m.buildRespWithPossibleChoicesAndTitle(possibleChoices)
	return resp, nil
}

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

func (m *Matcher) getAllPossibleChoices() map[int]struct{} {
	possibleChoices := make(map[int]struct{})
	for i := 0; i < m.choicesNum; i++ {
		possibleChoices[i] = struct{}{}
	}
	return possibleChoices
}

func (m *Matcher) buildRespWithPossibleChoicesAndTitle(choices map[int]struct{}) [][]string {
	resp := make([][]string, 0)
	resp = append(resp, m.originalDict[0])
	for k, _ := range choices {
		resp = append(resp, m.originalDict[k+1])
	}
	return resp
}

func (m *Matcher) getInvertedIndex() {

}

// separateQueries:
// queries like: C1 == "A" or C2 %26= "B"
// to
// [[C1,==,A,and][C2,&=,B,or]]
func separateQueries(queries string) ([][]string, error) {
	err := checkIfQueriesValid(queries)
	if err != nil {
		log.Printf("[seperateQueries] queries not invalid: %v", queries)
		return nil, errors.New("invalid query")
	}
	queryArr := make([][]string, 0)
	//todo: seperate
	return queryArr, nil
}

func checkIfQueriesValid(queries string) error {
	return nil
}

type Selector interface {
	selectWithQuery(*Matcher, []string) map[int]struct{}
}

type EqualSelector struct {
}
type NotEqualSelector struct {
}
type EqualInsensitiveCaseSelector struct {
}
type ContainSelector struct {
}

func (e *EqualSelector) selectWithQuery(m *Matcher, query []string) map[int]struct{} {
	possibleChoices := make(map[int]struct{})

	if query[0] != "*" {
		columns, ok := m.dict[query[0]]
		if !ok {
			log.Printf("[selectWithQuery] with title %v not found", query[0])
			return possibleChoices
		}
		for i := 0; i < m.choicesNum; i++ {
			if columns[i] == query[2] {
				possibleChoices[i] = struct{}{}
			}
		}
	} else {
		for i, choices := range m.originalDict[1:] {
			valid := true
			for _, choice := range choices {
				if choice != query[2] {
					valid = false
					break
				}
			}
			if valid {
				possibleChoices[i] = struct{}{}
			}
		}
	}
	return possibleChoices
}

func (e *NotEqualSelector) selectWithQuery(m *Matcher, query []string) map[int]struct{} {
	possibleChoices := make(map[int]struct{})

	if query[0] != "*" {
		columns, ok := m.dict[query[0]]
		if !ok {
			log.Printf("[selectWithQuery] with title %v not found", query[0])
			return possibleChoices
		}
		for i := 0; i < m.choicesNum; i++ {
			if columns[i] != query[2] {
				possibleChoices[i] = struct{}{}
			}
		}
	} else {
		for i, choices := range m.originalDict[1:] {
			valid := true
			for _, choice := range choices {
				if choice == query[2] {
					valid = false
					break
				}
			}
			if valid {
				possibleChoices[i] = struct{}{}
			}
		}
	}
	return possibleChoices
}

func (e *ContainSelector) selectWithQuery(m *Matcher, query []string) map[int]struct{} {
	possibleChoices := make(map[int]struct{})

	if query[0] != "*" {
		columns, ok := m.dict[query[0]]
		if !ok {
			log.Printf("[selectWithQuery] with title %v not found", query[0])
			return possibleChoices
		}
		for i := 0; i < m.choicesNum; i++ {
			if strings.Contains(columns[i], query[2]) {
				possibleChoices[i] = struct{}{}
			}
		}
	} else {
		for i, choices := range m.originalDict[1:] {
			valid := true
			for _, choice := range choices {
				if !strings.Contains(choice, query[2]) {
					valid = false
					break
				}
			}
			if valid {
				possibleChoices[i] = struct{}{}
			}
		}
	}
	return possibleChoices
}

func (e *EqualInsensitiveCaseSelector) selectWithQuery(m *Matcher, query []string) map[int]struct{} {
	possibleChoices := make(map[int]struct{})

	if query[0] != "*" {
		columns, ok := m.dict[query[0]]
		if !ok {
			log.Printf("[selectWithQuery] with title %v not found", query[0])
			return possibleChoices
		}
		for i := 0; i < m.choicesNum; i++ {
			if strings.ToLower(columns[i]) == strings.ToLower(query[2]) {
				possibleChoices[i] = struct{}{}
			}
		}
	} else {
		for i, choices := range m.originalDict[1:] {
			valid := true
			for _, choice := range choices {
				if strings.ToLower(choice) != strings.ToLower(query[2]) {
					valid = false
					break
				}
			}
			if valid {
				possibleChoices[i] = struct{}{}
			}
		}
	}
	return possibleChoices
}
