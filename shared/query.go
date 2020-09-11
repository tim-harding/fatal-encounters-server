package shared

import (
	"fmt"
	"strings"
)

// TODO: don't include LIMIT OFFSET if it has no contents

// NewQuery creates an empty query
func NewQuery() Query {
	return Query{
		subqueries: []Clauser{},
	}
}

// Query is a SQL query
type Query struct {
	subqueries []Clauser
}

// String constructs the SQL query string
func (q Query) String() string {
	subqueries := make([]string, 0, len(q.subqueries))
	for _, subquery := range q.subqueries {
		subqueries = append(subqueries, subquery.Term())
	}
	query := strings.Join(subqueries, " ")
	for i := range q.subqueries {
		placeholder := fmt.Sprintf("$%d", i+1)
		query = strings.Replace(query, "?", placeholder, 1)
	}
	return query
}

// Parms is the parameters for the query
func (q Query) Parms() []interface{} {
	parms := make([]interface{}, 0)
	for _, query := range q.subqueries {
		for _, clause := range query.Parms() {
			parms = append(parms, clause)
		}
	}
	return parms
}

// AddClause adds a new clause to the query
func (q *Query) AddClause(s Clauser) {
	q.subqueries = append(q.subqueries, s)
}

type selectClause struct {
	table string
	rows  []string
}

// NewSelectClause creates a SELECT FROM clause
func NewSelectClause(table string, rows []string) Clauser {
	return selectClause{
		table: table,
		rows:  rows,
	}
}

// Term returns a SQL snippet
func (s selectClause) Term() string {
	rows := strings.Join(s.rows, ", ")
	return fmt.Sprintf("SELECT %s FROM %s", rows, s.table)
}

// Parms returns the SQL query placeholder contents
func (s selectClause) Parms() []interface{} {
	return []interface{}{}
}

// Clauser is an interface for SQL query WHERE terms
type Clauser interface {
	Term() string
	Parms() []interface{}
}

// Page is a SQL query pagination term
type Page struct {
	limit  int
	offset int
}

// Term returns a SQL snippet
func (p Page) Term() string {
	return "LIMIT ? OFFSET ?"
}

// Parms returns the SQL query placeholder contents
func (p Page) Parms() []interface{} {
	return []interface{}{
		p.limit,
		p.offset,
	}
}

// Combinator is an enumeration of SQL WHERE combinators
type Combinator int

const (
	// CombinatorAnd is a SQL AND keyword
	CombinatorAnd Combinator = iota
	// CombinatorOr is a SQL OR keyword
	CombinatorOr
)

// WhereClause is a SQL where clause
type WhereClause struct {
	combinator Combinator
	clauses    []Clauser
}

// NewWhereClause creates a new WHERE clause
func NewWhereClause(combinator Combinator) WhereClause {
	return WhereClause{
		combinator: combinator,
		clauses:    []Clauser{},
	}
}

// AddClause adds a new clause to the WHERE statement
func (w *WhereClause) AddClause(clause Clauser) {
	w.clauses = append(w.clauses, clause)
}

// Term returns a SQL snippet
func (w WhereClause) Term() string {
	if len(w.clauses) == 0 {
		return ""
	}
	subqueries := make([]string, 0, len(w.clauses))
	for _, subquery := range w.clauses {
		subqueries = append(subqueries, subquery.Term())
	}
	combinator := [...]string{"AND", "OR"}[w.combinator]
	combinator = fmt.Sprintf(" %s ", combinator)
	combined := strings.Join(subqueries, combinator)
	return fmt.Sprintf("WHERE %s", combined)
}

// Parms returns the SQL query placeholder contents
func (w WhereClause) Parms() []interface{} {
	parms := make([]interface{}, 0)
	for _, subquery := range w.clauses {
		for _, parm := range subquery.Parms() {
			parms = append(parms, parm)
		}
	}
	return parms
}

type textSearchClause struct {
	column string
	term   string
}

// NewTextSearchClause creates a case insensitive text search term
func NewTextSearchClause(column, term string) Clauser {
	return textSearchClause{
		column: column,
		term:   term,
	}
}

// Term returns a SQL snippet
func (s textSearchClause) Term() string {
	return fmt.Sprintf("%s ILIKE '%%' || ? || '%%'", s.column)
}

// Parms returns the SQL query placeholder contents
func (s textSearchClause) Parms() []interface{} {
	return []interface{}{
		s.term,
	}
}
