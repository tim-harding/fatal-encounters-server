package city

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tim-harding/fatal-encounters-server/shared"
)

type cityFilter struct {
	Name  string `json:"name"`
	State int    `json:"state"`
}

type cityRow struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State int    `json:"state"`
}

type cityResponse struct {
	Rows []cityRow `json:"rows"`
}

const (
	cityQueryBase           = "SELECT id, name, state FROM city"
	cityQueryConditionBegin = " WHERE "
	cityQueryConditionAnd   = " AND "
	cityQueryFilterState    = "state=$"
	cityQueryFilterName     = "name ILIKE '%' || $1 || '%'"
	cityQueryLimit          = " LIMIT 12"
)

// City responds to /city queries
func City(w http.ResponseWriter, r *http.Request) {
	var filter cityFilter
	filter.State = -1

	stateStrs, ok := r.URL.Query()["state"]
	if ok && len(stateStrs) == 1 {
		stateStr := stateStrs[0]
		stateInt, err := strconv.Atoi(stateStr)
		if err == nil {
			filter.State = stateInt
		}
	}

	filterStrs, ok := r.URL.Query()["filter"]
	if ok && len(filterStrs) == 1 {
		filter.Name = filterStrs[0]
	}

	closure := queryCityFilterNone
	filterByState := filter.State > -1
	filterByName := len(filter.Name) > 0
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

	rows, err := closure(filter)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		id    int
		name  string
		state int
	)
	response := cityResponse{make([]cityRow, 0)}
	for rows.Next() {
		err := rows.Scan(&id, &name, &state)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		row := cityRow{id, name, state}
		response.Rows = append(response.Rows, row)
	}
	err = rows.Err()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func queryCityFilterState(filter cityFilter) (*sql.Rows, error) {
	query := fmt.Sprintf(
		"%v%v%v%v%v",
		cityQueryBase,
		cityQueryConditionBegin,
		cityQueryFilterState,
		"1",
		cityQueryLimit,
	)
	return shared.Db.Query(query, filter.State)
}

func queryCityFilterName(filter cityFilter) (*sql.Rows, error) {
	query := fmt.Sprintf(
		"%v%v%v%v",
		cityQueryBase,
		cityQueryConditionBegin,
		cityQueryFilterName,
		cityQueryLimit,
	)
	return shared.Db.Query(query, filter.Name)
}

func queryCityFilterBoth(filter cityFilter) (*sql.Rows, error) {
	query := fmt.Sprintf(
		"%v%v%v%v%v%v%v",
		cityQueryBase,
		cityQueryConditionBegin,
		cityQueryFilterName,
		cityQueryConditionAnd,
		cityQueryFilterState,
		"2",
		cityQueryLimit,
	)
	return shared.Db.Query(query, filter.Name, filter.State)
}

func queryCityFilterNone(filter cityFilter) (*sql.Rows, error) {
	query := fmt.Sprintf(
		"%v%v",
		cityQueryBase,
		cityQueryLimit,
	)
	return shared.Db.Query(query)
}
