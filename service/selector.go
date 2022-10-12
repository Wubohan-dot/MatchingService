package service

import (
	"log"
	"strings"
)

// Selector
// interface of selector, its implementations handle different operators of query
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

func selectorInit() {
	operator2Selector = make(map[string]Selector)
	operator2Selector["=="] = new(EqualSelector)
	operator2Selector["!="] = new(NotEqualSelector)
	operator2Selector["&="] = new(ContainSelector)
	operator2Selector["$="] = new(EqualInsensitiveCaseSelector)
}

func (e *EqualSelector) selectWithQuery(m *Matcher, query []string) map[int]struct{} {
	possibleChoices := make(map[int]struct{})

	if query[1] != "*" {
		columns, ok := m.dict[query[1]]
		if !ok {
			log.Printf("[selectWithQuery] with title %v not found", query[1])
			return possibleChoices
		}
		for i := 0; i < m.choicesNum; i++ {
			if columns[i] == query[3] {
				possibleChoices[i] = struct{}{}
			}
		}
	} else {
		for i, choices := range m.originalDict[1:] {
			valid := true
			for _, choice := range choices {
				if choice != query[3] {
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

	if query[1] != "*" {
		columns, ok := m.dict[query[1]]
		if !ok {
			log.Printf("[selectWithQuery] with title %v not found", query[1])
			return possibleChoices
		}
		for i := 0; i < m.choicesNum; i++ {
			if columns[i] != query[3] {
				possibleChoices[i] = struct{}{}
			}
		}
	} else {
		for i, choices := range m.originalDict[1:] {
			valid := true
			for _, choice := range choices {
				if choice == query[3] {
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

	if query[1] != "*" {
		columns, ok := m.dict[query[1]]
		if !ok {
			log.Printf("[selectWithQuery] with title %v not found", query[1])
			return possibleChoices
		}
		for i := 0; i < m.choicesNum; i++ {
			if strings.Contains(columns[i], query[3]) {
				possibleChoices[i] = struct{}{}
			}
		}
	} else {
		for i, choices := range m.originalDict[1:] {
			valid := true
			for _, choice := range choices {
				if !strings.Contains(choice, query[3]) {
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

	if query[1] != "*" {
		columns, ok := m.dict[query[1]]
		if !ok {
			log.Printf("[selectWithQuery] with title %v not found", query[1])
			return possibleChoices
		}
		for i := 0; i < m.choicesNum; i++ {
			if strings.ToLower(columns[i]) == strings.ToLower(query[3]) {
				possibleChoices[i] = struct{}{}
			}
		}
	} else {
		for i, choices := range m.originalDict[1:] {
			valid := true
			for _, choice := range choices {
				if strings.ToLower(choice) != strings.ToLower(query[3]) {
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
