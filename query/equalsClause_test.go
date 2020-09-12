package query

import "testing"

func TestEqualsClause(t *testing.T) {
	query := baseQuery()
	w := NewWhereClause(CombinatorAnd)
	equals := NewEqualsClause("column", 3)
	w.AddClause(equals)
	query.AddClause(w)
	const wanted = "SELECT a, b FROM test WHERE column = $1"
	try(query, wanted, t)
}
