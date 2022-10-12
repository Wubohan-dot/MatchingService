package service

import (
	"fmt"
	"testing"
)

func TestReadDict(t *testing.T) {
}

func TestCheckIfQueriesValid(t *testing.T) {
	queryArr, err := checkIfQueriesValid("C1 == \"A\"")
	fmt.Println(err)
	fmt.Println(queryArr)
	queryArr, err = checkIfQueriesValid("C1 == \"A\" wbh")
	fmt.Println(err)
	fmt.Println(queryArr)
}

// “==” equal, “!=” not equal, “$=” equal (case insensitive), “&=” contain (the query term is a substring of the data cell).
func TestMatchWithQueries(t *testing.T) {
	matcher.originalDict = [][]string{{"A", "B", "C"}, {"a1", "b1", "c1"}, {"a2", "b2", "c2"}}
	matcher.dict = map[string][]string{
		"A": {"a1", "a2"},
		"B": {"b1", "b2"},
		"C": {"c1", "c2"},
	}
	matcher.choicesNum = 2
	matcher.columnNum = 3

	// test with simple single condition
	// should all return resp:[[A B C] [a1 b1 c1]]	err:<nil>
	successCases := []string{
		"C == \"c1\"", "C != \"c2\"", "C $= \"C1\"", "C &= \"1\"",
	}
	for _, successCase := range successCases {
		resp, err := matcher.MatchWithQueries(successCase)
		fmt.Printf("case:%v\tresp:%v\terr:%v\n", successCase, resp, err)
	}

	// test with simple *
	// should all return resp:[[A B C] [a2 b2 c2]]	err:<nil>
	successCases = []string{
		"* &= \"2\"",
	}
	for _, successCase := range successCases {
		resp, err := matcher.MatchWithQueries(successCase)
		fmt.Printf("case:%v\tresp:%v\terr:%v\n", successCase, resp, err)
	}

	// test with simple *, but no data should be found
	// should all return resp:[[A B C]]	err:<nil>
	successCases = []string{
		"* &= \"3\"",
	}
	for _, successCase := range successCases {
		resp, err := matcher.MatchWithQueries(successCase)
		fmt.Printf("case:%v\tresp:%v\terr:%v\n", successCase, resp, err)
	}

	// test with "and"
	// resp:[[A B C] [a2 b2 c2]]	err:<nil>
	// resp:[[A B C]]	err:<nil>
	successCases = []string{
		"C == \"c2\" and B == \"b2\"", "C == \"c2\" and B == \"b1\"",
	}
	for _, successCase := range successCases {
		resp, err := matcher.MatchWithQueries(successCase)
		fmt.Printf("case:%v\tresp:%v\terr:%v\n", successCase, resp, err)
	}

	// test with "or"
	// resp:[[A B C] [a2 b2 c2]]	err:<nil>
	// resp:[[A B C] [a2 b2 c2] [a1 b1 c1]]	err:<nil>
	successCases = []string{
		"C == \"c2\" or B == \"www\"", "C == \"c2\" or B == \"b1\"",
	}
	for _, successCase := range successCases {
		resp, err := matcher.MatchWithQueries(successCase)
		fmt.Printf("case:%v\tresp:%v\terr:%v\n", successCase, resp, err)
	}

	// test with both "or" and "and"
	// resp:[[A B C]]	err:<nil>
	// resp:[[A B C] [a1 b1 c1]]	err:<nil>
	successCases = []string{
		"C == \"c2\" or B == \"www\" and B == \"www\"", "C == \"c2\" or B == \"b1\" and A == \"a1\"",
	}
	for _, successCase := range successCases {
		resp, err := matcher.MatchWithQueries(successCase)
		fmt.Printf("case:%v\tresp:%v\terr:%v\n", successCase, resp, err)
	}
}
