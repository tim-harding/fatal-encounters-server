package incidentroute

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type coordinate struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"long"`
}

type mappingData struct {
	ID         int        `json:"id"`
	Coordinate coordinate `json:"coordinate"`
}

// ID only (for filtering map pins)
// Basic info (for listed results)
// Full info (for details page)

// HandleRouteMapping delivers basic incident information
func HandleRouteMapping(w http.ResponseWriter, r *http.Request) {
	shared.HandleRoute(w, r, buildQuery, translateRow)
}

func buildQuery(r *http.Request) query.Clauser {
	q := query.NewQuery()
	q.AddClause(selectClause())
	q.AddClause(whereClause(r))
	q.AddClause(shared.LimitClause(r))
	return q
}

func selectClause() query.Clauser {
	desiredRowNames := []string{
		"id",
		"latitude",
		"longitude",
	}
	return query.NewSelectClause("incident", desiredRowNames)
}

func whereClause(r *http.Request) query.Clauser {
	w := query.NewWhereClause(query.CombinatorAnd)
	w.AddClause(stateClause(r))
	w.AddClause(shared.SearchClause(r))
	return w
}

func stateClause(r *http.Request) query.Clauser {
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
	return query.NewInClause("state", states)
}

func translateRow(rows *sql.Rows) (interface{}, error) {
	var (
		id        int
		latitude  float32
		longitude float32
	)
	err := rows.Scan(&id, &latitude, &longitude)
	if err != nil {
		return nil, err
	}
	coord := coordinate{latitude, longitude}
	return mappingData{id, coord}, nil
}
