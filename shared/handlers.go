package shared

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// QueryBuilderFunc creates a query for the HTTP request
type QueryBuilderFunc func(r *http.Request) Clauser

// RowTranslatorFunc creates a row of JSON response from a database row
type RowTranslatorFunc func(rows *sql.Rows) (interface{}, error)

type response struct {
	Rows []interface{} `json:"rows"`
}

func newResponse() response {
	return response{[]interface{}{}}
}

// HandleRoute responds to queries
func HandleRoute(w http.ResponseWriter, r *http.Request, queryBuilder QueryBuilderFunc, rowTranslator RowTranslatorFunc) {
	res, err := buildResponse(r, queryBuilder, rowTranslator)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func buildResponse(r *http.Request, queryBuilder QueryBuilderFunc, rowTranslator RowTranslatorFunc) (interface{}, error) {
	query := queryBuilder(r)
	queryString := query.String()
	log.Printf("Database query: %s", queryString)

	rows, err := Db.Query(queryString, query.Parameters()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res, err := rowsToResponse(rows, rowTranslator)
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func rowsToResponse(rows *sql.Rows, rowTranslator RowTranslatorFunc) (interface{}, error) {
	res := newResponse()
	for rows.Next() {
		row, err := rowTranslator(rows)
		if err != nil {
			return nil, err
		}
		res.Rows = append(res.Rows, row)
	}
	return res, nil
}
