package shared

import "testing"

func TestTextSearchQuery(t *testing.T) {
	search := NewTextSearchClause("column", "find")
	where := NewWhereClause(CombinatorAnd)
	where.AddClause(search)
	query := baseQuery()
	query.AddClause(where)
	const wanted = "SELECT a, b FROM test WHERE column ILIKE '%' || $1 || '%'"
	try(query, wanted, t)
}

func TestIgnoreTextSearchIfSearchTermIsEmpty(t *testing.T) {
	search := NewTextSearchClause("column", "")
	where := NewWhereClause(CombinatorAnd)
	where.AddClause(search)
	query := baseQuery()
	query.AddClause(where)
	const wanted = "SELECT a, b FROM test"
	try(query, wanted, t)
}
