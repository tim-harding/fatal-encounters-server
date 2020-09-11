package shared

import "testing"

func base() selectClause {
	return selectClause{
		table: "test",
		rows:  []string{"a", "b"},
	}
}

func baseQuery() Query {
	query := NewQuery()
	query.AddClause(base())
	return query
}

func try(query Query, wanted string, t *testing.T) {
	sql := query.String()
	if sql != wanted {
		t.Errorf("Was `%s`;\nWant `%s`", sql, wanted)
	}
}

func TestBuildsBasicQuery(t *testing.T) {
	query := baseQuery()
	const wanted = "SELECT a, b FROM test"
	try(query, wanted, t)
}

func TestPagination(t *testing.T) {
	page := NewPageClause(12, 1)
	query := baseQuery()
	query.AddClause(page)
	const wanted = "SELECT a, b FROM test LIMIT $1 OFFSET $2"
	try(query, wanted, t)
}

func TestTextSearchQuery(t *testing.T) {
	search := NewTextSearchClause("column", "find")
	where := NewWhereClause(CombinatorAnd)
	where.AddClause(search)
	query := baseQuery()
	query.AddClause(where)
	const wanted = "SELECT a, b FROM test WHERE column ILIKE '%' || $1 || '%'"
	try(query, wanted, t)
}

func TestEmptyWhereClauseDoesNothing(t *testing.T) {
	where := NewWhereClause(CombinatorAnd)
	query := baseQuery()
	query.AddClause(where)
	// Todo: trim query?
	const wanted = "SELECT a, b FROM test "
	try(query, wanted, t)
}
