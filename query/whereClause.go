package query

import (
	"fmt"
)

// Combinator is an enumeration of SQL WHERE combinators
type Combinator int

const (
	// CombinatorAnd is a SQL AND keyword
	CombinatorAnd Combinator = iota
	// CombinatorOr is a SQL OR keyword
	CombinatorOr
)

type whereClause struct {
	expr subexpression
}

// NewWhereClause creates a new WHERE clause
func NewWhereClause(combinator Combinator) Subclauser {
	tokens := []string{
		" AND ",
		" OR ",
	}
	combinatorToken := tokens[combinator]
	expr := newSubexpression(combinatorToken)
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
