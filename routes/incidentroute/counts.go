package incidentroute

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type countFor struct {
	Key   int `json:"key"`
	Count int `json:"count"`
}

type countsResponse struct {
	Counts map[string][]countFor `json:"counts"`
	Rows   []int                 `json:"rows"`
}

// HandleCountRoute handles requests to /incident/count
func HandleCountRoute(w http.ResponseWriter, r *http.Request) {
	tx, err := shared.Db.Begin()
	if err != nil {
		shared.InternalError(w, err)
		return
	}
	_, err = tx.Exec(sqlDropTemp)
	if err != nil {
		shared.InternalError(w, err)
		return
	}
	_, err = tx.Exec(sqlCreateTemp)
	if err != nil {
		shared.InternalError(w, err)
		return
	}
	q := populateFiltered(r)
	_, err = tx.Exec(q.String(), q.Parameters()...)
	if err != nil {
		shared.InternalError(w, err)
		return
	}
	ids, err := queryIds(tx)
	if err != nil {
		shared.InternalError(w, err)
		return
	}
	counts := map[string][]countFor{}
	for _, order := range orderColumns {
		q := countQuery(order)
		count, err := queryCountFor(q, tx)
		if err != nil {
			shared.InternalError(w, err)
			return
		}
		counts[order.Name] = count
	}
	tx.Commit()
	res := countsResponse{counts, ids}
	json.NewEncoder(w).Encode(res)
}

func queryIds(tx *sql.Tx) ([]int, error) {
	str := allFiltered().String()
	log.Printf(str)
	rows, err := tx.Query(str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []int{}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func queryCountFor(query query.Clauser, tx *sql.Tx) ([]countFor, error) {
	str := query.String()
	log.Printf(str)
	rows, err := tx.Query(str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []countFor{}
	for rows.Next() {
		count := countFor{}
		err := rows.Scan(&count.Key, &count.Count)
		if err != nil {
			return nil, err
		}
		out = append(out, count)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func populateFiltered(r *http.Request) query.Clauser {
	q := query.NewQuery()
	q.AddClause(query.NewInsertClause("filtered"))
	q.AddClause(query.NewSelectClause("incident", []string{"incident.id"}))
	q.AddClause(whereClauseFilter(r))
	q.AddClause(orderClause(r))
	return q
}

func allFiltered() query.Clauser {
	q := query.NewQuery()
	q.AddClause(query.NewSelectClause("filtered", []string{"id"}))
	return q
}

func countQuery(column orderColumn) query.Clauser {
	q := query.NewQuery()
	columns := []string{
		column.Translated(),
		"COUNT(1)",
	}
	q.AddClause(query.NewSelectClause("incident", columns))
	filtered := fmt.Sprintf(sqlFiltered, column.Column)
	q.AddClause(query.NewRawSQL(filtered))
	return q
}

const (
	sqlDropTemp   = "DROP TABLE IF EXISTS filtered"
	sqlCreateTemp = `
		CREATE TEMPORARY TABLE IF NOT EXISTS filtered (
			id INTEGER PRIMARY KEY NOT NULL
		)
		ON COMMIT DROP
	`
	sqlFiltered = `
		WHERE incident.id
		IN (
			SELECT filtered.id
			FROM filtered
		)
		AND %s IS NOT NULL
		GROUP BY 1
		ORDER BY 1
	`
)

type orderColumn struct {
	Name       string
	Column     string
	translator string
}

func (o *orderColumn) Translated() string {
	return fmt.Sprintf(o.translator, o.Column)
}

var (
	orderColumns = [...]orderColumn{
		{
			"race",
			"race_id",
			"%s",
		},
		{
			"cause",
			"cause_id",
			"%s",
		},
		{
			"year",
			"date",
			"EXTRACT(YEAR FROM incident.%s) AS yyyy",
		},
		{
			"age",
			"age",
			"%s",
		},
	}
)
