package incidentroute

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type enum struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type maybeEnum struct {
	ID   *int
	Name *string
}

type detailRow struct {
	ID          int       `json:"id"`
	Name        *string   `json:"name"`
	Age         *int      `json:"age"`
	Date        time.Time `json:"date"`
	ImageURL    *string   `json:"imageUrl"`
	IsMale      *bool     `json:"isMale"`
	Address     *string   `json:"address"`
	Description string    `json:"description"`
	ArticleURL  *string   `json:"articleUrl"`
	VideoURL    *string   `json:"videoUrl"`
	Zipcode     *int      `json:"zipcode"`
	Cause       enum      `json:"cause"`
	UseOfForce  enum      `json:"useOfForce"`
	Race        *enum     `json:"race"`
	County      *enum     `json:"county"`
	Agency      *enum     `json:"agency"`
	City        *enum     `json:"city"`
}

// HandleIncidentDetailRoute responds to /incident/{id} routes
func HandleIncidentDetailRoute(w http.ResponseWriter, r *http.Request) {
	shared.HandleIDRoute(w, r, buildDetailQuery(r), translateDetailRow, "incident")
}

func buildDetailQuery(r *http.Request) query.Clauser {
	q := query.NewSubexpression(" ")
	q.AddClause(selectClause(rowKindDetail))
	q.AddClause(joinClausesDetail())
	return q
}

func joinClausesDetail() query.Clauser {
	expr := query.NewSubexpression(" ")
	for _, table := range enumTables {
		clause := query.NewJoinClause(table)
		expr.AddClause(clause)
	}
	return expr
}

func translateDetailRow(rows *sql.Rows) (interface{}, error) {
	row := detailRow{}

	enums := make([]maybeEnum, 4)
	targets := []**enum{
		&row.Race,
		&row.County,
		&row.Agency,
		&row.City,
	}

	err := rows.Scan(
		&row.ID,
		&row.Name,
		&row.Age,
		&row.Date,
		&row.ImageURL,
		&row.IsMale,
		&row.Address,
		&row.Description,
		&row.ArticleURL,
		&row.VideoURL,
		&row.Zipcode,

		&row.Cause.ID,
		&row.Cause.Name,

		&row.UseOfForce.ID,
		&row.UseOfForce.Name,

		&enums[0].ID,
		&enums[0].Name,

		&enums[1].ID,
		&enums[1].Name,

		&enums[2].ID,
		&enums[2].Name,

		&enums[3].ID,
		&enums[3].Name,
	)

	if err != nil {
		return nil, err
	}

	for i, maybe := range enums {
		if maybe.ID != nil {
			*targets[i] = &enum{
				*maybe.ID,
				*maybe.Name,
			}
		}
	}

	return row, err
}
