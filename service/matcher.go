package service

import (
	"encoding/csv"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"
)

// Matcher
// the structure responsible for matching task
// originalDict: to store the original data in CSV file, for example, [[A,B,C],[a1,b1,c1],[a2,b2,c2]]
// dict: to store a map from title to its column, to help with processing tasks, for example [A:[a1,a2],B:[b1,b2],C:[c1,c2]]
// columnNum indicates number of columns, since A,B,C, it's 3
// choicesNum indicates rows in the table(except for title row), it's 2
type Matcher struct {
	originalDict  [][]string
	dict          map[string][]string
	columnNum     int
	choicesNum    int
	invertedIndex map[string][]string
}

var matcher Matcher
var operator2Selector map[string]Selector

// init
// to init service package, mainly including a matcher and a map of selector
// the matcher.dict and the
func init() {
	matcher.dict = make(map[string][]string)
	matcher.getDictAndOriginalDict("./dict.csv")
	matcher.getInvertedIndex()
	selectorInit()
	rand.Seed(time.Now().Unix())
}

// MatcherInstance
// to return an existing matcher in service package
func MatcherInstance() *Matcher {
	return &matcher
}

// getDictAndOriginalDict
// read from file and fill in matcher.dict and matcher.originalDict
func (m *Matcher) getDictAndOriginalDict(fileName string) error {
	records, err := m.readFile(fileName)
	if err != nil {
		log.Fatalf("[getDictAndOriginalDict] err:%v", err)
		return err
	}
	m.originalDict = records
	err = m.formatToDict(records)
	if err != nil {
		log.Fatalf("[getDictAndOriginalDict] err:%v", err)
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

// checkIfRecordsValid
// check if the data read from CSV file valid or not, it requires:
// 1. One row for title + at least one row for data
// 2. The column name(title) only be characters or digits (A-Z, a-z, 0-9, case sensitive)
// 3. Each row has the same length
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

// MatchWithQueries
// the main function to handle task
// accepts queries in a whole string and returns answer rows with column name
func (m *Matcher) MatchWithQueries(queries string) ([][]string, error) {
	queryArr, err := m.separateQueries(queries)
	if err != nil {
		log.Printf("[MatchWithQueries] cannot separate queries to array: %v\n", queries)
		return nil, err
	}
	// queryNum := len(queryArr)
	possibleChoices := m.getAllPossibleChoices()
	//  “==” equal, “!=” not equal, “$=” equal (case insensitive), “&=” contain (the query term is a substring of the data cell)
	for _, query := range queryArr {
		selector, ok := operator2Selector[query[2]]
		if !ok {
			return nil, errors.New("wrong operator")
		}
		tmpPossibleChoices := selector.selectWithQuery(&matcher, query)
		if query[0] == "and" {
			possibleChoices = getIntersection(tmpPossibleChoices, possibleChoices)
		} else {
			possibleChoices = getUnion(tmpPossibleChoices, possibleChoices)
		}

	}
	resp := m.buildRespWithPossibleChoicesAndTitle(possibleChoices)
	return resp, nil
}

// separateQueries:
// queries like: C1 == "A" or C2 %26= "B"
// to
// [[and,C1,==,A][or,C2,&=,B]]
func (m *Matcher) separateQueries(queries string) ([][]string, error) {
	words, err := checkIfQueriesValid(queries)
	if err != nil {
		log.Printf("[seperateQueries] queries not invalid: %v", queries)
		return nil, err
	}
	queryArr := make([][]string, 0)
	for i := 0; i < len(words); i += 4 {
		queryArr = append(queryArr, words[i:i+4])
	}
	return queryArr, nil
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
