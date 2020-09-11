package cityroute

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/tim-harding/fatal-encounters-server/shared"
)

// Todo: add LIMIT OFFSET

type city struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State int    `json:"state"`
}

type response struct {
	Rows []city `json:"rows"`
}

type stateClause struct {
	state int
}

func (s stateClause) Term() string {
	return "state = ?"
}

func (s stateClause) Parms() []interface{} {
	return []interface{}{s.state}
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

func buildQuery(r *http.Request) shared.Query {
	w := shared.NewWhereClause(shared.CombinatorAnd)
	addStateClause(r, &w)
	addSearchClause(r, &w)
	q := shared.NewQuery()
	rows := []string{
		"id",
		"name",
		"state",
	}
	base := shared.NewSelectClause("city", rows)
	q.AddClause(base)
	q.AddClause(w)
	return q
}

func addStateClause(r *http.Request, w *shared.WhereClause) {
	strings, ok := r.URL.Query()["state"]
	if ok && len(strings) == 1 {
		string := strings[0]
		integer, err := strconv.Atoi(string)
		if err == nil {
			clause := stateClause{integer}
			w.AddClause(clause)
		}
	}
}

func addSearchClause(r *http.Request, w *shared.WhereClause) {
	strings, ok := r.URL.Query()["search"]
	if ok && len(strings) == 1 {
		clause := shared.NewTextSearchClause("name", strings[0])
		w.AddClause(clause)
	}
}

func responseForQuery(query shared.Query) (response, error) {
	queryString := query.String()
	log.Printf(queryString)
	rows, err := shared.Db.Query(queryString, query.Parms()...)
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
