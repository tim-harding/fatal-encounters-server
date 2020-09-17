package query

import (
	"fmt"
	"strings"
)

// Clauser is an interface for SQL query WHERE terms
type Clauser interface {
	fmt.Stringer
	Parameters() []interface{}
}

// Subclauser is a clause that can have clauses added to it
type Subclauser interface {
	Clauser
	AddClause(clause Clauser)
}

// NewQuery creates an empty query
func NewQuery() Subclauser {
	return &query{newSubexpression(" ")}
}

type query struct {
	expr Subclauser
}

func (q *query) String() string {
	query := q.expr.String()
	for i := range q.Parameters() {
		placeholder := fmt.Sprintf("$%d", i+1)
		query = strings.Replace(query, "?", placeholder, 1)
	}
	return query
}

func (q *query) Parameters() []interface{} {
	return q.expr.Parameters()
}

// AddClause adds a new clause to the query
func (q *query) AddClause(s Clauser) {
	q.expr.AddClause(s)
}
