package incidentroute

import (
	"database/sql"
	"net/http"

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
	enumTables := []string{
		"agency",
		"cause",
		"city",
		"county",
		"race",
		"state",
		"use_of_force",
	}
	for _, table := range enumTables {
		clause := shared.InClause(r, table)
		w.AddClause(clause)
	}
	w.AddClause(shared.SearchClause(r))
	w.AddClause(ageClause(r, "ageMin", query.ComparatorGreaterEqual))
	w.AddClause(ageClause(r, "ageMax", query.ComparatorLesserEqual))
	w.AddClause(genderMaskClause(r))
	return w
}

func ageClause(r *http.Request, key string, comparator query.Comparator) query.Clauser {
	ok, value := shared.MaybeQueryInt(r, key)
	if !ok {
		return nil
	}
	return query.NewCompareClause(comparator, "age", value)
}

func genderMaskClause(r *http.Request) query.Clauser {
	querystrings, ok := r.URL.Query()["gender"]
	if ok && len(querystrings) > 0 {
		var male bool
		switch querystrings[0] {
		case "male":
			male = true
			break
		case "female":
			male = false
			break
		default:
			return nil
		}
		return query.NewCompareClause(query.ComparatorEqual, "is_male", male)
	}
	return nil
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
