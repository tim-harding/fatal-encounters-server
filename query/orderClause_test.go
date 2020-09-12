package query

import "testing"

func TestOrderClause(t *testing.T) {
	query := baseQuery()
	order := NewOrderClause(OrderingDescending, []string{"id", "thing"})
	query.AddClause(order)
	const wanted = "SELECT a, b FROM test ORDER BY id, thing DESC"
	try(query, wanted, t)
}
