package shared

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
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
func HandleRoute(w http.ResponseWriter, r *http.Request, query query.Clauser, rowTranslator RowTranslatorFunc) {
	res, err := buildResponse(r, query, rowTranslator)
	if err != nil {
		InternalError(w, err)
		return
	}
	json.NewEncoder(w).Encode(res)
}

// InternalError sends an internal server error message
func InternalError(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusInternalServerError)
}

// Error sends an error response
func Error(w http.ResponseWriter, err error, code int) {
	log.Printf("%v", err)
	message := http.StatusText(code)
	http.Error(w, message, code)
}

func buildResponse(r *http.Request, query query.Clauser, rowTranslator RowTranslatorFunc) (interface{}, error) {
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
	querystrings, ok := r.URL.Query()[key]
	if !ok || len(querystrings) < 1 {
		return defaultValue
	}
	value, err := strconv.Atoi(querystrings[0])
	if err != nil {
		return defaultValue
	}
	return value
}

// MaybeQueryInt gets an integer value if available
func MaybeQueryInt(r *http.Request, key string) (bool, int) {
	querystrings, ok := r.URL.Query()[key]
	if !ok || len(querystrings) < 1 {
		return false, -1
	}
	integer, err := strconv.Atoi(querystrings[0])
	if err != nil {
		return false, -1
	}
	return true, integer
}

// SearchClause creates a text search clause based on the `name` column
func SearchClause(r *http.Request) query.Clauser {
	strings, ok := r.URL.Query()["search"]
	if !ok || len(strings) < 1 {
		return nil
	}
	return query.NewTextSearchClause("name", strings[0])
}

// InClause creates an IN clause from the request
func InClause(r *http.Request, column string) query.Clauser {
	values := queryInts(r, column)
	return query.NewInClause(column, values)
}

func queryInts(r *http.Request, key string) []int {
	mask := make([]int, 0)
	querystrings, ok := r.URL.Query()[key]
	if ok {
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
	return mask
}

// HandleIDRoute creates a handler function for ID routes
func HandleIDRoute(w http.ResponseWriter, r *http.Request, selectClause query.Clauser, rowTranslator RowTranslatorFunc) {
	query, err := buildWhereQuery(selectClause, r)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}
	HandleRoute(w, r, query, rowTranslator)
}

func buildWhereQuery(base query.Clauser, r *http.Request) (query.Clauser, error) {
	w, err := whereClauseID(r)
	if err != nil {
		return nil, err
	}
	q := query.NewQuery()
	q.AddClause(base)
	q.AddClause(w)
	return q, nil
}

func whereClauseID(r *http.Request) (query.Clauser, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	match := query.NewCompareClause(query.ComparisonEqual, "id", id)
	w := query.NewWhereClause(query.CombinatorAnd)
	w.AddClause(match)
	return w, nil
}

// IgnoreClause sets up the query to reject certain IDs from the response
func IgnoreClause(r *http.Request) query.Clauser {
	values := queryInts(r, "ignore")
	in := query.NewInClause("id", values)
	not := query.NewNotClause(in)
	return not
}
