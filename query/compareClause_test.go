package query

import "testing"

func TestCompareClause(t *testing.T) {
	query := baseQuery()
	w := NewWhereClause(CombinatorAnd)
	equals := NewCompareClause(ComparatorEqual, "column", 3)
	w.AddClause(equals)
	query.AddClause(w)
	const wanted = "SELECT a, b FROM test WHERE column = $1"
	try(query, wanted, t)
}
