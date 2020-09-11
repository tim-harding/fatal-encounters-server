package shared

import "testing"

func TestBuildsBasicQuery(t *testing.T) {
	base := Select{
		table: "test",
		rows:  []string{"a", "b"},
	}
	query := RawQuery{
		subqueries: []Subquerier{
			base,
		},
	}
	const wanted = "SELECT a, b FROM test"
	sql := query.Build()
	if sql != wanted {
		t.Errorf("Was `%s`;\nWant `%s`", sql, wanted)
	}
}
