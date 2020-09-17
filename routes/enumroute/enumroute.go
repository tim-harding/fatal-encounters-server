package enumroute

import (
	"database/sql"
	"net/http"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type enum struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// HandleBaseRouteFactory creates functions to respond to queries
// on enumeration tables that include id and name
func HandleBaseRouteFactory(tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := buildQuery(r, tableName)
		shared.HandleRoute(w, r, query, translateRow)
	}
}

// HandleIDRouteFactory creates functions to respond to queries
// on enumeration tables that include id and name
func HandleIDRouteFactory(table string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := buildQuery(r, table)
		shared.HandleIDRoute(w, r, query, translateRow, "tableName")
	}
}

func buildQuery(r *http.Request, table string) query.Clauser {
	q := query.NewQuery()
	q.AddClause(selectClause(table))
	q.AddClause(whereClause(r, table))
	q.AddClause(orderClause())
	q.AddClause(shared.LimitClause(r))
	return q
}

func selectClause(tableName string) query.Clauser {
	desiredRowNames := []string{
		"id",
		"name",
	}
	return query.NewSelectClause(tableName, desiredRowNames)
}

func whereClause(r *http.Request, table string) query.Clauser {
	w := query.NewWhereClause(query.CombinatorAnd)
	w.AddClause(shared.SearchClause(r))
	w.AddClause(shared.IgnoreClause(r, table))
	return w
}

func orderClause() query.Clauser {
	order := query.OrderingAscending
	columns := []string{"name"}
	return query.NewOrderClause(order, columns)
}

func translateRow(rows *sql.Rows) (interface{}, error) {
	var (
		id   int
		name string
	)
	err := rows.Scan(&id, &name)
	if err != nil {
		return nil, err
	}
	row := enum{id, name}
	return row, nil
}
