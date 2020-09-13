package query

import (
	"fmt"
)

type compareClause struct {
	comparison Comparison
	column     string
	match      interface{}
}

// Comparison is an enumeration of comparison options
type Comparison int

const (
	// ComparisonEqual compares values using `=` operator
	ComparisonEqual Comparison = iota
	// ComparisonGreater compares values using `>` operator
	ComparisonGreater
	// ComparisonLesser compares values using `<` operator
	ComparisonLesser
	// ComparisonGreaterEqual compares values using `>=` operator
	ComparisonGreaterEqual
	// ComparisonLesserEqual compares values using `>=` operator
	ComparisonLesserEqual
)

var comparatorStrings = []string{
	"=",
	">",
	"<",
	">=",
	"<=",
}

// NewCompareClause creates a `column = ?` SQL clause
func NewCompareClause(comparison Comparison, column string, match interface{}) Clauser {
	return &compareClause{comparison, column, match}
}

func (c *compareClause) String() string {
	comparisonString := comparatorStrings[c.comparison]
	return fmt.Sprintf("%s %s ?", c.column, comparisonString)
}

func (c *compareClause) Parameters() []interface{} {
	return []interface{}{c.match}
}
