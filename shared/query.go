package shared

import (
	"fmt"
	"strings"
)

// Todo: limit=-1 for all results

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
		subqueries = append(subqueries, subquery.String())
	}
	query := strings.Join(subqueries, " ")
	for i := range q.subqueries {
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

type pageClause struct {
	limit  int
	offset int
}

// NewPageClause creates a pagination clause
func NewPageClause(limit, offset int) Clauser {
	return &pageClause{limit, offset}
}

// String returns a SQL snippet
func (p *pageClause) String() string {
	if p.offset > 0 {
		return "LIMIT ? OFFSET ?"
	}
	return "LIMIT ?"
}

// Parameters returns the SQL query placeholder contents
func (p *pageClause) Parameters() []interface{} {
	if p.offset > 0 {
		return []interface{}{
			p.limit,
			p.offset,
		}
	}
	return []interface{}{
		p.limit,
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
	if len(w.clauses) == 0 {
		return ""
	}
	subqueries := make([]string, 0, len(w.clauses))
	for _, subquery := range w.clauses {
		subqueries = append(subqueries, subquery.String())
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

type textSearchClause struct {
	column string
	term   string
}

// NewTextSearchClause creates a case insensitive text search term
func NewTextSearchClause(column, term string) Clauser {
	return &textSearchClause{
		column: column,
		term:   term,
	}
}

// String returns a SQL snippet
func (s *textSearchClause) String() string {
	return fmt.Sprintf("%s ILIKE '%%' || ? || '%%'", s.column)
}

// Parameters returns the SQL query placeholder contents
func (s *textSearchClause) Parameters() []interface{} {
	return []interface{}{
		s.term,
	}
}
