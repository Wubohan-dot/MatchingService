package main

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
