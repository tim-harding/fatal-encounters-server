package cityRoute

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tim-harding/fatal-encounters-server/shared"
)

type queryFunction func(filter filter) (*sql.Rows, error)

type filter struct {
	Search string `json:"name"`
	State  int    `json:"state"`
}

type city struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State int    `json:"state"`
}

type response struct {
	Rows []city `json:"rows"`
}

var (
	queryByState        *sql.Stmt
	queryByName         *sql.Stmt
	queryByNameAndState *sql.Stmt
	queryAny            *sql.Stmt
)

// HandleRoute responds to /city queries
func HandleRoute(w http.ResponseWriter, r *http.Request) {
	filter := createFilter(r)
	query := pickQueryFunction(filter)
	res, err := responseForQuery(query, filter)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func init() {
	const (
		queryComponentBase           = "SELECT id, name, state FROM city"
		queryComponentConditionBegin = " WHERE "
		queryComponentConditionAnd   = " AND "
		queryComponentFilterState    = "state=$"
		queryComponentFilterName     = "name ILIKE '%' || $1 || '%'"
		queryComponentLimit          = " LIMIT 12"
	)

	query := fmt.Sprintf(
		"%v%v%v%v%v",
		queryComponentBase,
		queryComponentConditionBegin,
		queryComponentFilterState,
		"1",
		queryComponentLimit,
	)
	queryByState = statementForQuery(query)

	query = fmt.Sprintf(
		"%v%v%v%v",
		queryComponentBase,
		queryComponentConditionBegin,
		queryComponentFilterName,
		queryComponentLimit,
	)
	queryByName = statementForQuery(query)

	query = fmt.Sprintf(
		"%v%v%v%v%v%v%v",
		queryComponentBase,
		queryComponentConditionBegin,
		queryComponentFilterName,
		queryComponentConditionAnd,
		queryComponentFilterState,
		"2",
		queryComponentLimit,
	)
	queryByNameAndState = statementForQuery(query)

	query = fmt.Sprintf(
		"%v%v",
		queryComponentBase,
		queryComponentLimit,
	)
	queryAny = statementForQuery(query)
}

func statementForQuery(query string) *sql.Stmt {
	stmt, err := shared.Db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

func responseForQuery(query queryFunction, filter filter) (response, error) {
	rows, err := query(filter)
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

func pickQueryFunction(filter filter) queryFunction {
	closure := queryCityFilterNone
	filterByState := filter.State > -1
	filterByName := len(filter.Search) > 0
	if filterByState {
		if filterByName {
			closure = queryCityFilterBoth
		} else {
			closure = queryCityFilterState
		}
	} else {
		if filterByName {
			closure = queryCityFilterName
		}
	}
	return closure
}

func createFilter(r *http.Request) filter {
	filter := filter{"", -1}
	populateStateFilter(r, &filter)
	populateSearchFilter(r, &filter)
	return filter
}

func populateStateFilter(r *http.Request, filter *filter) {
	strings, ok := r.URL.Query()["state"]
	if ok && len(strings) == 1 {
		string := strings[0]
		integer, err := strconv.Atoi(string)
		if err == nil {
			filter.State = integer
		}
	}
}

func populateSearchFilter(r *http.Request, filter *filter) {
	strings, ok := r.URL.Query()["search"]
	if ok && len(strings) == 1 {
		filter.Search = strings[0]
	}
}

func queryCityFilterState(filter filter) (*sql.Rows, error) {
	return queryByState.Query(filter.State)
}

func queryCityFilterName(filter filter) (*sql.Rows, error) {
	return queryByName.Query(filter.Search)
}

func queryCityFilterBoth(filter filter) (*sql.Rows, error) {
	return queryByNameAndState.Query(filter.Search, filter.State)
}

func queryCityFilterNone(filter filter) (*sql.Rows, error) {
	return queryAny.Query()
}
