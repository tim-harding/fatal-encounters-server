package query

import "testing"

func TestEmptyWhereClauseDoesNothing(t *testing.T) {
	where := NewWhereClause(CombinatorAnd)
	query := baseQuery()
	query.AddClause(where)
	const wanted = "SELECT a, b FROM test"
	try(query, wanted, t)
}

func TestCombinesMultipleWhereClauses(t *testing.T) {
	where := NewWhereClause(CombinatorAnd)
	textSearch := NewTextSearchClause("column", "things")
	in := NewInClause("column", []int{3, 5, 7})
	where.AddClause(textSearch)
	where.AddClause(in)
	query := baseQuery()
	query.AddClause(where)
	const wanted = "SELECT a, b FROM test " +
		"WHERE column ILIKE '%' || $1 || '%' AND " +
		"column IN ($2, $3, $4)"
	try(query, wanted, t)
}
