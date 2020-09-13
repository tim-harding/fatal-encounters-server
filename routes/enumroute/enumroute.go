package enumroute

import (
	"database/sql"
	"net/http"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type state struct {
	ID   int
	Name string
}

// HandleRouteFactory creates functions to respond to queries
// on enumeration tables that include id and name
func HandleRouteFactory(tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := buildQuery(r, tableName)
		shared.HandleRoute(w, r, query, translateRow)
	}
}

func buildQuery(r *http.Request, tableName string) query.Clauser {
	q := query.NewQuery()
	q.AddClause(selectClause(tableName))
	q.AddClause(whereClause(r))
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

func whereClause(r *http.Request) query.Clauser {
	w := query.NewWhereClause(query.CombinatorAnd)
	w.AddClause(shared.SearchClause(r))
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
	row := state{id, name}
	return row, nil
}
