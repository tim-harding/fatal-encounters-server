package shared

import (
	"fmt"
	"strings"
)

// Todo: struct to cache query permutations as prepared statements
// Todo: replace `?` with `$1`, `$2`, etc

// RawQuery is a generator for SQL query strings
type RawQuery struct {
	subqueries []Subquerier
}

// Build constructs the SQL query string
func (q *RawQuery) Build() string {
	subqueries := make([]string, 0, len(q.subqueries))
	for _, subquery := range q.subqueries {
		subqueries = append(subqueries, subquery.Subquery())
	}
	query := strings.Join(subqueries, " ")
	for i := range q.subqueries {
		placeholder := fmt.Sprintf("$%d", i)
		query = strings.Replace(query, "?", placeholder, 1)
	}
	return query
}

// Select is the SQL SELECT FROM statement
type Select struct {
	table string
	rows  []string
}

// Subquery returns a SQL snippet
func (s Select) Subquery() string {
	rows := strings.Join(s.rows, ", ")
	return fmt.Sprintf("SELECT %s FROM %s", rows, s.table)
}

// Parms returns the SQL query placeholder contents
func (s Select) Parms() []interface{} {
	return []interface{}{}
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
func (p Page) Subquery() string {
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

// Where is a SQL where clause
type Where struct {
	combinator Combinator
	subqueries []Subquerier
}

// Subquery returns a SQL snippet
func (w Where) Subquery() string {
	subqueries := make([]string, 0, len(w.subqueries))
	for _, subquery := range w.subqueries {
		subqueries = append(subqueries, subquery.Subquery())
	}
	combinator := [...]string{"AND", "OR"}[w.combinator]
	combinator = fmt.Sprintf(" %s ", combinator)
	combined := strings.Join(subqueries, combinator)
	return fmt.Sprintf("(%s)", combined)
}

// Parms returns the SQL query placeholder contents
func (w Where) Parms() []interface{} {
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
func (s SearchFilter) Subquery() string {
	return fmt.Sprintf("%s ILIKE '%%' || ? || '%%'", s.column)
}

// Parms returns the SQL query placeholder contents
func (s SearchFilter) Parms() []interface{} {
	return []interface{}{
		s.term,
	}
}
