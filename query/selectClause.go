package query

import (
	"fmt"
	"strings"
)

type selectClause struct {
	table string
	rows  []string
}

// NewSelectClause creates a SELECT FROM clause
func NewSelectClause(table string, rows []string) Clauser {
	return &selectClause{table, rows}
}

// String returns a SQL snippet
func (s *selectClause) String() string {
	rows := strings.Join(s.rows, ", ")
	return fmt.Sprintf("SELECT %s FROM %s", rows, s.table)
}

// Parameters returns the SQL query placeholder contents
func (s *selectClause) Parameters() []interface{} {
	return []interface{}{}
}
