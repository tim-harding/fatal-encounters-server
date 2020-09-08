package enumroute

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tim-harding/fatal-encounters-server/shared"
)

type state struct {
	ID   int
	Name string
}

type response struct {
	Rows []state
}

// HandleRouteFactory creates functions to respond to queries
// on enumeration tables that include id and name
func HandleRouteFactory(tableName string) http.HandlerFunc {
	query := createStatement(tableName)
	return func(w http.ResponseWriter, r *http.Request) {
		handleRoute(w, r, query)
	}
}

func handleRoute(w http.ResponseWriter, r *http.Request, query *sql.Stmt) {
	res, err := responseForQuery(query)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func createStatement(tableName string) *sql.Stmt {
	queryString := fmt.Sprintf("SELECT id, name FROM %s", tableName)
	stmt, err := shared.Db.Prepare(queryString)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

func responseForQuery(query *sql.Stmt) (response, error) {
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
