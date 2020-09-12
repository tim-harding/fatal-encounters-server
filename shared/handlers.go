package shared

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tim-harding/fatal-encounters-server/query"
)

// QueryBuilderFunc creates a query for the HTTP request
type QueryBuilderFunc func(r *http.Request) query.Clauser

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

// LimitClause creates a limit clause from the request
func LimitClause(r *http.Request) query.Clauser {
	limit := QueryInt(r, "count", 6)
	page := QueryInt(r, "page", 0)
	offset := page * limit
	clause := query.NewPageClause(limit, offset)
	return clause
}

// QueryInt gets an integer value from the request query string
func QueryInt(r *http.Request, key string, defaultValue int) int {
	value := defaultValue
	querystrings, ok := r.URL.Query()[key]
	if ok && len(querystrings) > 0 {
		querystring := querystrings[0]
		integer, err := strconv.Atoi(querystring)
		if err == nil {
			value = integer
		}
	}
	return value
}

// MaybeQueryInt gets an integer value if available
func MaybeQueryInt(r *http.Request, key string) (bool, int) {
	querystrings, ok := r.URL.Query()[key]
	if ok && len(querystrings) > 0 {
		querystring := querystrings[0]
		integer, err := strconv.Atoi(querystring)
		if err == nil {
			return true, integer
		}
	}
	return false, -1
}

// SearchClause creates a text search clause based on the `name` column
func SearchClause(r *http.Request) query.Clauser {
	strings, ok := r.URL.Query()["search"]
	if ok && len(strings) > 0 {
		return query.NewTextSearchClause("name", strings[0])
	}
	return nil
}

// InClause creates an IN clause from the request
func InClause(r *http.Request, column string) query.Clauser {
	mask := make([]int, 0)
	querystrings, ok := r.URL.Query()[column]
	if ok && len(querystrings) > 0 {
		for _, querystring := range querystrings {
			parts := strings.Split(querystring, ",")
			for _, part := range parts {
				integer, err := strconv.Atoi(part)
				if err == nil {
					mask = append(mask, integer)
				}
			}
		}
	}
	return query.NewInClause(column, mask)
}
