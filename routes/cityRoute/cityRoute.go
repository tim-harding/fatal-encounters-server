package cityroute

import (
	"database/sql"
	"net/http"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type city struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State int    `json:"state"`
}

// HandleRoute responds to /city queries
func HandleRoute(w http.ResponseWriter, r *http.Request) {
	shared.HandleRoute(w, r, buildQuery, translateRow)
}

func buildQuery(r *http.Request) query.Clauser {
	q := query.NewQuery()
	q.AddClause(selectClause())
	q.AddClause(whereClause(r))
	q.AddClause(orderClause())
	q.AddClause(shared.LimitClause(r))
	return q
}

func selectClause() query.Clauser {
	desiredRowNames := []string{
		"id",
		"name",
		"state",
	}
	return query.NewSelectClause("city", desiredRowNames)
}

func whereClause(r *http.Request) query.Clauser {
	w := query.NewWhereClause(query.CombinatorAnd)
	w.AddClause(shared.InClause(r, "state"))
	w.AddClause(shared.SearchClause(r))
	return w
}

func orderClause() query.Clauser {
	order := query.OrderingAscending
	columns := []string{"name", "state"}
	return query.NewOrderClause(order, columns)
}

func translateRow(rows *sql.Rows) (interface{}, error) {
	var (
		id    int
		name  string
		state int
	)
	err := rows.Scan(&id, &name, &state)
	if err != nil {
		return nil, err
	}
	return city{id, name, state}, nil
}
