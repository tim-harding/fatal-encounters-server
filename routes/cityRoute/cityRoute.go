package cityroute

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/tim-harding/fatal-encounters-server/shared"
)

type city struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State int    `json:"state"`
}

type response struct {
	Rows []city `json:"rows"`
}

// HandleRoute responds to /city queries
func HandleRoute(w http.ResponseWriter, r *http.Request) {
	query := buildQuery(r)
	res, err := responseForQuery(query)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func buildQuery(r *http.Request) shared.Clauser {
	rows := []string{
		"id",
		"name",
		"state",
	}
	q := shared.NewQuery()
	base := shared.NewSelectClause("city", rows)
	q.AddClause(base)
	w := shared.NewWhereClause(shared.CombinatorAnd)
	addClause(w, r, makeStateClause)
	addClause(w, r, makeSearchClause)
	q.AddClause(w)
	addClause(q, r, makeLimitClause)
	return q
}

type maybeClause func(r *http.Request) shared.Clauser

func addClause(w shared.Subclauser, r *http.Request, clauseFunc maybeClause) {
	clause := clauseFunc(r)
	if clause != nil {
		w.AddClause(clause)
	}
}

func makeLimitClause(r *http.Request) shared.Clauser {
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

func makeStateClause(r *http.Request) shared.Clauser {
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

func makeSearchClause(r *http.Request) shared.Clauser {
	strings, ok := r.URL.Query()["search"]
	if ok && len(strings) == 1 {
		return shared.NewTextSearchClause("name", strings[0])
	}
	return nil
}

func responseForQuery(query shared.Clauser) (response, error) {
	queryString := query.String()
	log.Printf("Database query: %s", queryString)
	rows, err := shared.Db.Query(queryString, query.Parameters()...)
	if err != nil {
		return response{}, err
	}

	defer rows.Close()

	res, err := rowsToResponse(rows)
	if err != nil {
		return response{}, err
	}

	err = rows.Err()
	if err != nil {
		return response{}, err
	}

	return res, nil
}

func rowsToResponse(rows *sql.Rows) (response, error) {
	res := response{make([]city, 0)}
	for rows.Next() {
		var (
			id    int
			name  string
			state int
		)
		err := rows.Scan(&id, &name, &state)
		if err != nil {
			return response{}, err
		}
		row := city{id, name, state}
		res.Rows = append(res.Rows, row)
	}
	return res, nil
}
