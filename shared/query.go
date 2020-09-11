package shared

import (
	"fmt"
	"strings"
)

// Todo: struct to cache query permutations as prepared statements
// Todo: replace `?` with `$1`, `$2`, etc

// QueryMaker is a generator for SQL query strings
type QueryMaker struct {
	table      string
	rows       []string
	subqueries []Subquerier
}

// Build constructs the SQL query string
func (q *QueryMaker) Build() string {
	rows := strings.Join(q.rows, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s", rows, q.table)
	for i := range q.subqueries {
		placeholder := fmt.Sprintf("$%d", i)
		query = strings.Replace(query, "?", placeholder, 1)
	}
	return query
}

// Subquerier is an interface for SQL query WHERE terms
type Subquerier interface {
	Subquery() string
	Parms() []interface{}
}

// Page is a SQL query pagination term
type Page struct {
	limit  int
	offset int
}

// Subquery returns a SQL snippet
func (p *Page) Subquery() string {
	return "LIMIT ? OFFSET ?"
}

// Parms returns the SQL query placeholder contents
func (p *Page) Parms() []interface{} {
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

// Where is a SQL where clause
type Where struct {
	combinator Combinator
	subqueries []Subquerier
}

// Subquery returns a SQL snippet
func (w *Where) Subquery() string {
	combinator := [...]string{"AND", "OR"}[w.combinator]
	subqueries := make([]string, len(w.subqueries))
	for _, subquery := range w.subqueries {
		subqueries = append(subqueries, subquery.Subquery())
	}
	combined := strings.Join(subqueries, combinator)
	return fmt.Sprintf("(%s)", combined)
}

// Parms returns the SQL query placeholder contents
func (w *Where) Parms() []interface{} {
	parms := make([]interface{}, 0)
	for _, subquery := range w.subqueries {
		for _, parm := range subquery.Parms() {
			parms = append(parms, parm)
		}
	}
	return parms
}

// SearchFilter is a case insensitive text search filter
type SearchFilter struct {
	column string
	term   string
}

// Subquery returns a SQL snippet
func (s *SearchFilter) Subquery() string {
	return fmt.Sprintf("%s ILIKE '%%' || ? || '%%'", s.column)
}

// Parms returns the SQL query placeholder contents
func (s *SearchFilter) Parms() []interface{} {
	return []interface{}{
		s.term,
	}
}
