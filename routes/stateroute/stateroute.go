package stateroute

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type state struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Shortname string `json:"shortname"`
}

var desiredRowNames = [...]string{
	"id",
	"name",
	"shortname",
}

// HandleBaseRoute responds to /city queries
func HandleBaseRoute(w http.ResponseWriter, r *http.Request) {
	shared.HandleRoute(w, r, buildQuery(r), translateRow)
}

// HandleIDRoute responds to /city queries
func HandleIDRoute(w http.ResponseWriter, r *http.Request) {
	shared.HandleIDRoute(w, r, selectClause(), translateRow, "state")
}

func buildQuery(r *http.Request) query.Clauser {
	q := query.NewQuery()
	q.AddClause(selectClause())
	q.AddClause(whereClause(r))
	q.AddClause(orderClause())
	q.AddClause(shared.LimitClause(r))
	return q
}

func selectClause() query.Clauser {
	return query.NewSelectClause("state", desiredRowNames[:])
}

func whereClause(r *http.Request) query.Clauser {
	w := query.NewWhereClause(query.CombinatorAnd)
	or := query.NewConditionsClause(query.CombinatorOr)
	for _, clause := range searchClauses(r) {
		or.AddClause(clause)
	}
	w.AddClause(or)
	w.AddClause(shared.IgnoreClause(r, "state"))
	return w
}

func orderClause() query.Clauser {
	order := query.OrderingAscending
	columns := []string{"name"}
	return query.NewOrderClause(order, columns)
}

func translateRow(rows *sql.Rows) (interface{}, error) {
	var (
		id        int
		name      string
		shortname []byte
	)
	err := rows.Scan(&id, &name, &shortname)
	if err != nil {
		return nil, err
	}
	shortnameStr := string(shortname)
	return state{id, name, shortnameStr}, nil
}

func searchClauses(r *http.Request) []query.Clauser {
	querystringValues, ok := r.URL.Query()["search"]
	if ok && len(querystringValues) == 1 {
		term := querystringValues[0]
		nameSearch := query.NewTextSearchClause("name", term)
		clauses := []query.Clauser{nameSearch}
		if len(term) == 2 {
			caps := strings.ToUpper(term)
			shortnameSearch := query.NewCompareClause(query.ComparisonEqual, "shortname", caps)
			clauses = append(clauses, shortnameSearch)
		}
		return clauses
	}
	return []query.Clauser{}
}
