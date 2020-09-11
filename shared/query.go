package shared

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
	return &query{
		subqueries: []Clauser{},
	}
}

// Query is a SQL query
type query struct {
	subqueries []Clauser
}

// String constructs the SQL query string
func (q *query) String() string {
	subqueries := make([]string, 0, len(q.subqueries))
	for _, subquery := range q.subqueries {
		querystring := subquery.String()
		if querystring != "" {
			subqueries = append(subqueries, querystring)
		}
	}
	query := strings.Join(subqueries, " ")
	for i := range q.Parameters() {
		placeholder := fmt.Sprintf("$%d", i+1)
		query = strings.Replace(query, "?", placeholder, 1)
	}
	return query
}

// Parameters is the parameters for the query
func (q *query) Parameters() []interface{} {
	parameters := make([]interface{}, 0)
	for _, query := range q.subqueries {
		for _, clause := range query.Parameters() {
			parameters = append(parameters, clause)
		}
	}
	return parameters
}

// AddClause adds a new clause to the query
func (q *query) AddClause(s Clauser) {
	q.subqueries = append(q.subqueries, s)
}
