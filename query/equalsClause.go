package query

import (
	"fmt"
)

type equalsClause struct {
	column string
	match  interface{}
}

// NewEqualsClause creates a `column = ?` SQL clause
func NewEqualsClause(column string, match interface{}) Clauser {
	return &equalsClause{column, match}
}

func (m *equalsClause) String() string {
	return fmt.Sprintf("%s = ?", m.column)
}

func (m *equalsClause) Parameters() []interface{} {
	return []interface{}{m.match}
}
