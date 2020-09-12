package query

import "testing"

func TestInClause(t *testing.T) {
	query := baseQuery()
	where := NewWhereClause(CombinatorAnd)
	in := NewInClause("column", []int{3, 5, 7})
	where.AddClause(in)
	query.AddClause(where)
	const wanted = "SELECT a, b FROM test WHERE column IN ($1, $2, $3)"
	try(query, wanted, t)
}

func TestIgnoreInClause(t *testing.T) {
	query := baseQuery()
	where := NewWhereClause(CombinatorAnd)
	in := NewInClause("column", []int{})
	where.AddClause(in)
	query.AddClause(where)
	const wanted = "SELECT a, b FROM test"
	try(query, wanted, t)
}
