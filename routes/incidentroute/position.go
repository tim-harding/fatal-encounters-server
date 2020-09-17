package incidentroute

import (
	"database/sql"
	"net/http"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type position struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

type positionRow struct {
	ID       int      `json:"id"`
	Position position `json:"position"`
}

// HandleIncidentPositionRoute handles requests to /incident/position
func HandleIncidentPositionRoute(w http.ResponseWriter, r *http.Request) {
	shared.HandleRoute(w, r, buildPositionQuery(), translatePositionRow)
}

func buildPositionQuery() query.Clauser {
	q := query.NewQuery()
	q.AddClause(selectClause(rowKindPosition))
	return q
}

func translatePositionRow(rows *sql.Rows) (interface{}, error) {
	row := positionRow{}
	err := rows.Scan(
		&row.ID,
		&row.Position.Latitude,
		&row.Position.Longitude,
	)
	return row, err
}
