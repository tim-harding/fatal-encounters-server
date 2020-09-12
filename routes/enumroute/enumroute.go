package enumroute

import (
	"database/sql"
	"net/http"

	"github.com/tim-harding/fatal-encounters-server/shared"
)

type state struct {
	ID   int
	Name string
}

// HandleRouteFactory creates functions to respond to queries
// on enumeration tables that include id and name
func HandleRouteFactory(tableName string) http.HandlerFunc {
	queryBuilder := queryBuilderFactory(tableName)
	return func(w http.ResponseWriter, r *http.Request) {
		shared.HandleRoute(w, r, queryBuilder, translateRow)
	}
}

func queryBuilderFactory(tableName string) shared.QueryBuilderFunc {
	return func(r *http.Request) shared.Clauser {
		q := shared.NewQuery()
		q.AddClause(selectClause(tableName))
		q.AddClause(whereClause(r))
		q.AddClause(orderClause())
		q.AddClause(shared.LimitClause(r))
		return q
	}
}

func selectClause(tableName string) shared.Clauser {
	desiredRowNames := []string{
		"id",
		"name",
	}
	return shared.NewSelectClause(tableName, desiredRowNames)
}

func whereClause(r *http.Request) shared.Clauser {
	w := shared.NewWhereClause(shared.CombinatorAnd)
	w.AddClause(shared.SearchClause(r))
	return w
}

func orderClause() shared.Clauser {
	order := shared.OrderingAscending
	columns := []string{"name"}
	return shared.NewOrderClause(order, columns)
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
