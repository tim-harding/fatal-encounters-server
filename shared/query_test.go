package shared

import "testing"

func base() Select {
	return Select{
		table: "test",
		rows:  []string{"a", "b"},
	}
}

func try(query RawQuery, wanted string, t *testing.T) {
	sql := query.Build()
	if sql != wanted {
		t.Errorf("Was `%s`;\nWant `%s`", sql, wanted)
	}
}

func TestBuildsBasicQuery(t *testing.T) {
	query := RawQuery{
		subqueries: []Subquerier{
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
	query := RawQuery{
		subqueries: []Subquerier{
			base(),
			page,
		},
	}
	const wanted = "SELECT a, b FROM test LIMIT $1 OFFSET $2"
	try(query, wanted, t)
}

func TestTextSearchQuery(t *testing.T) {
	search := SearchFilter{
		column: "column",
		term:   "find",
	}
	where := Where{
		combinator: CombinatorAnd,
		subqueries: []Subquerier{
			search,
		},
	}
	query := RawQuery{
		subqueries: []Subquerier{
			base(),
			where,
		},
	}
	const wanted = "SELECT a, b FROM test WHERE column ILIKE '%' || $1 || '%'"
	try(query, wanted, t)
}
