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
	q.AddClause(limitClause(r))
	q.AddClause(orderClause())
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
	w.AddClause(searchClause(r))
	return w
}

func limitClause(r *http.Request) shared.Clauser {
	limit := 1
	strings, ok := r.URL.Query()["count"]
	if ok && len(strings) == 1 {
		string := strings[0]
		integer, err := strconv.Atoi(string)
		if err == nil {
			limit = integer
		}
	}
	clause := shared.NewPageClause(limit, 0)
	return clause
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

func searchClause(r *http.Request) shared.Clauser {
	strings, ok := r.URL.Query()["search"]
	if ok && len(strings) == 1 {
		return shared.NewTextSearchClause("name", strings[0])
	}
	return nil
}

func orderClause() shared.Clauser {
	order := shared.OrderingAscending
	columns := []string{"state", "id"}
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
