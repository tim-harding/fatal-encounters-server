package incidentroute

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

// HandleIncidentFilterRoute responds to /incident/{id} routes
func HandleIncidentFilterRoute(w http.ResponseWriter, r *http.Request) {
	query := buildFilterQuery(r)
	shared.HandleRoute(w, r, query, translateFilterRow)
}

func buildFilterQuery(r *http.Request) query.Clauser {
	q := query.NewQuery()
	q.AddClause(selectClause(rowKindFilter))
	// TODO: Only include join if filtering by state
	q.AddClause(query.NewJoinClause("city"))
	q.AddClause(whereClauseFilter(r))
	q.AddClause(orderClause(r))
	return q
}

func whereClauseFilter(r *http.Request) query.Clauser {
	w := query.NewWhereClause(query.CombinatorAnd)
	for _, table := range idQueryTables {
		column := fmt.Sprintf("%s_id", table)
		clause := shared.InClause(r, column)
		w.AddClause(clause)
	}
	w.AddClause(shared.SearchClause(r))
	w.AddClause(ageClause(r, "ageMin", query.ComparisonGreaterEqual))
	w.AddClause(ageClause(r, "ageMax", query.ComparisonLesserEqual))
	w.AddClause(genderMaskClause(r))
	w.AddClause(dateMaskClause(r, "dateMin", query.ComparisonGreaterEqual))
	w.AddClause(dateMaskClause(r, "dateMax", query.ComparisonLesserEqual))
	return w
}

func ageClause(r *http.Request, key string, comparator query.Comparison) query.Clauser {
	ok, value := shared.MaybeQueryInt(r, key)
	if !ok {
		return nil
	}
	return query.NewCompareClause(comparator, "age", value)
}

func genderMaskClause(r *http.Request) query.Clauser {
	querystrings, ok := r.URL.Query()["gender"]
	if !ok {
		return nil
	}
	isMale, ok := genders[querystrings[0]]
	if !ok {
		return nil
	}
	return query.NewCompareClause(query.ComparisonEqual, "is_male", isMale)
}

func orderClause(r *http.Request) query.Clauser {
	kind := pickOrderKind(r)
	column := orderKindColumns[kind]
	column = fmt.Sprintf("incident.%s", column)
	direction := pickOrderDirection(r)
	return query.NewOrderClause(direction, []string{column})
}

func pickOrderKind(r *http.Request) orderKind {
	querystrings, ok := r.URL.Query()["order"]
	if !ok {
		return orderKindID
	}
	order, ok := querystringToOrderKind[querystrings[0]]
	if !ok {
		return orderKindID
	}
	return order
}

func pickOrderDirection(r *http.Request) query.Ordering {
	querystrings, ok := r.URL.Query()["orderDirection"]
	if !ok {
		return query.OrderingAscending
	}
	orderDirection, ok := querystringToOrderDirection[querystrings[0]]
	if !ok {
		return query.OrderingAscending
	}
	return orderDirection
}

func dateMaskClause(r *http.Request, key string, comparator query.Comparison) query.Clauser {
	querystrings, ok := r.URL.Query()[key]
	if !ok {
		return nil
	}
	t, err := time.Parse("2006-Jan-02", querystrings[0])
	if err != nil {
		return nil
	}
	return query.NewCompareClause(comparator, "date", t)
}

func translateFilterRow(rows *sql.Rows) (interface{}, error) {
	var id int
	err := rows.Scan(
		&id,
	)
	return id, err
}
