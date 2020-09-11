package shared

import "testing"

func TestPagination(t *testing.T) {
	page := NewPageClause(12, 1)
	query := baseQuery()
	query.AddClause(page)
	const wanted = "SELECT a, b FROM test LIMIT $1 OFFSET $2"
	try(query, wanted, t)
}

func TestIgnoresOffsetWhenZero(t *testing.T) {
	page := NewPageClause(12, 0)
	query := baseQuery()
	query.AddClause(page)
	const wanted = "SELECT a, b FROM test LIMIT $1"
	try(query, wanted, t)
}

func TestIgnoresPaginationWhenLimitIsZero(t *testing.T) {
	page := NewPageClause(0, 0)
	query := baseQuery()
	query.AddClause(page)
	const wanted = "SELECT a, b FROM test"
	try(query, wanted, t)
}
