package cityroute

import (
	"database/sql"
	"net/http"
	"strconv"

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

func buildQuery(r *http.Request) shared.Clauser {
	q := shared.NewQuery()
	q.AddClause(selectClause())
	q.AddClause(whereClause(r))
	q.AddClause(orderClause())
	q.AddClause(shared.LimitClause(r))
	return q
}

func selectClause() shared.Clauser {
	desiredRowNames := []string{
		"id",
		"name",
		"state",
	}
	return shared.NewSelectClause("city", desiredRowNames)
}

func whereClause(r *http.Request) shared.Clauser {
	w := shared.NewWhereClause(shared.CombinatorAnd)
	w.AddClause(stateClause(r))
	w.AddClause(shared.SearchClause(r))
	return w
}

func stateClause(r *http.Request) shared.Clauser {
	states := make([]int, 0)
	strings, ok := r.URL.Query()["state"]
	if ok {
		for _, string := range strings {
			integer, err := strconv.Atoi(string)
			if err == nil {
				states = append(states, integer)
			}
		}
	}
	return shared.NewInClause("state", states)
}

func orderClause() shared.Clauser {
	order := shared.OrderingAscending
	columns := []string{"name", "state"}
	return shared.NewOrderClause(order, columns)
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
