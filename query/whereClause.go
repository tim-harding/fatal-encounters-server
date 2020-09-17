package query

import (
	"fmt"
)

type whereClause struct {
	expr Subclauser
}

// NewWhereClause creates a new WHERE clause
func NewWhereClause(combinator Combinator) Subclauser {
	expr := NewConditionsClause(combinator)
	return &whereClause{expr}
}

// AddClause adds a new clause to the WHERE statement
func (w *whereClause) AddClause(clause Clauser) {
	w.expr.AddClause(clause)
}

// String returns a SQL snippet
func (w *whereClause) String() string {
	expr := w.expr.String()
	if expr == "" {
		return ""
	}
	return fmt.Sprintf("WHERE %s", expr)
}

// Parameters returns the SQL query placeholder contents
func (w *whereClause) Parameters() []interface{} {
	return w.expr.Parameters()
}
