package stateRoute

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/tim-harding/fatal-encounters-server/shared"
)

var query *sql.Stmt

type state struct {
	ID   int
	Name string
}

type response struct {
	Rows []state
}

// HandleRoute responds to /city queries
func HandleRoute(w http.ResponseWriter, r *http.Request) {
	res, err := responseForQuery()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func init() {
	stmt, err := shared.Db.Prepare("SELECT id, name FROM state")
	if err != nil {
		log.Fatal(err)
	}
	query = stmt
}

func responseForQuery() (response, error) {
	rows, err := query.Query()
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
	res := response{make([]state, 0)}
	for rows.Next() {
		var (
			id   int
			name string
		)
		err := rows.Scan(&id, &name)
		if err != nil {
			return response{}, err
		}
		row := state{id, name}
		res.Rows = append(res.Rows, row)
	}
	return res, nil
}
