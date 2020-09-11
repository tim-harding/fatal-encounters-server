package shared

import "testing"

func base() selectClause {
	return selectClause{
		table: "test",
		rows:  []string{"a", "b"},
	}
}

func try(query Query, wanted string, t *testing.T) {
	sql := query.String()
	if sql != wanted {
		t.Errorf("Was `%s`;\nWant `%s`", sql, wanted)
	}
}

func TestBuildsBasicQuery(t *testing.T) {
	query := Query{
		subqueries: []Clauser{
			base(),
		},
	}
	const wanted = "SELECT a, b FROM test"
	try(query, wanted, t)
}

func TestPagination(t *testing.T) {
	page := Page{
		limit:  12,
		offset: 1,
	}
	query := Query{
		subqueries: []Clauser{
			base(),
			page,
		},
	}
	const wanted = "SELECT a, b FROM test LIMIT $1 OFFSET $2"
	try(query, wanted, t)
}

func TestTextSearchQuery(t *testing.T) {
	search := textSearchClause{
		column: "column",
		term:   "find",
	}
	where := WhereClause{
		combinator: CombinatorAnd,
		clauses: []Clauser{
			search,
		},
	}
	query := Query{
		subqueries: []Clauser{
			base(),
			where,
		},
	}
	const wanted = "SELECT a, b FROM test WHERE column ILIKE '%' || $1 || '%'"
	try(query, wanted, t)
}

func TestEmptyWhereClauseDoesNothing(t *testing.T) {
	where := NewWhereClause(CombinatorAnd)
	query := Query{
		subqueries: []Clauser{
			base(),
			where,
		},
	}
	// Todo: trim query?
	const wanted = "SELECT a, b FROM test "
	try(query, wanted, t)
}
