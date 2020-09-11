package shared

import "fmt"

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
