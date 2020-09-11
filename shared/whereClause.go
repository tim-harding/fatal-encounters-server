package shared

import (
	"fmt"
	"strings"
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
	combinator Combinator
	clauses    []Clauser
}

// NewWhereClause creates a new WHERE clause
func NewWhereClause(combinator Combinator) Subclauser {
	return &whereClause{
		combinator: combinator,
		clauses:    []Clauser{},
	}
}

// AddClause adds a new clause to the WHERE statement
func (w *whereClause) AddClause(clause Clauser) {
	w.clauses = append(w.clauses, clause)
}

// String returns a SQL snippet
func (w *whereClause) String() string {
	subqueries := make([]string, 0, len(w.clauses))
	for _, subquery := range w.clauses {
		querystring := subquery.String()
		if querystring != "" {
			subqueries = append(subqueries, querystring)
		}
	}
	if len(subqueries) < 1 {
		return ""
	}
	combinator := [...]string{"AND", "OR"}[w.combinator]
	combinator = fmt.Sprintf(" %s ", combinator)
	combined := strings.Join(subqueries, combinator)
	return fmt.Sprintf("WHERE %s", combined)
}

// Parameters returns the SQL query placeholder contents
func (w *whereClause) Parameters() []interface{} {
	parameters := make([]interface{}, 0)
	for _, subquery := range w.clauses {
		for _, parm := range subquery.Parameters() {
			parameters = append(parameters, parm)
		}
	}
	return parameters
}
