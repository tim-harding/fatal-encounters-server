package shared

import "testing"

func TestBuildsBasicQuery(t *testing.T) {
	query := QueryMaker{
		table:      "test",
		rows:       []string{"a", "b"},
		subqueries: []Subquerier{},
	}
	const wanted = "SELECT a, b FROM test"
	sql := query.Build()
	if sql != wanted {
		t.Errorf("Was `%s`;\nWant `%s`", sql, wanted)
	}
}
