package query

import (
	"fmt"
)

type compareClause struct {
	comparator Comparator
	column     string
	match      interface{}
}

// Comparator is an enumeration of comparison options
type Comparator int

const (
	// ComparatorEqual compares values using `=` operator
	ComparatorEqual Comparator = iota
	// ComparatorGreater compares values using `>` operator
	ComparatorGreater
	// ComparatorLesser compares values using `<` operator
	ComparatorLesser
	// ComparatorGreaterEqual compares values using `>=` operator
	ComparatorGreaterEqual
	// ComparatorLesserEqual compares values using `>=` operator
	ComparatorLesserEqual
)

var comparatorStrings = []string{
	"=",
	">",
	"<",
	">=",
	"<=",
}

// NewCompareClause creates a `column = ?` SQL clause
func NewCompareClause(comparator Comparator, column string, match interface{}) Clauser {
	return &compareClause{comparator, column, match}
}

func (m *compareClause) String() string {
	comparatorString := comparatorStrings[m.comparator]
	return fmt.Sprintf("%s %s ?", m.column, comparatorString)
}

func (m *compareClause) Parameters() []interface{} {
	return []interface{}{m.match}
}
