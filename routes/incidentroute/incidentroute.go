package incidentroute

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type position struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

type idRow struct {
	ID int `json:"id"`
}

type positionRow struct {
	idRow
	Position position `json:"position"`
}

type enum struct {
	idRow
	Name string `json:"name"`
}

type maybeEnum struct {
	ID   *int
	Name *string
}

type detailRow struct {
	idRow
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

type responseRow interface {
	FromRow(rows *sql.Rows)
}

// HandleRouteBase responds to /incident/ routes
func HandleRouteBase(w http.ResponseWriter, r *http.Request) {
	kind := pickRowKind(r)
	query := buildBaseQuery(r, kind)
	translateRow := translateFunc(kind)
	shared.HandleRoute(w, r, query, translateRow)
}

// HandleRouteID responds to /incident/{id} routes
func HandleRouteID(w http.ResponseWriter, r *http.Request) {
	kind := pickRowKind(r)
	selectClause := selectClause(kind)
	translateRow := translateFunc(kind)
	shared.HandleIDRoute(w, r, selectClause, translateRow)
}

func pickRowKind(r *http.Request) rowKind {
	querystrings, ok := r.URL.Query()["rowKind"]
	if !ok {
		return rowKindID
	}
	kind, ok := querystringToRowKinds[querystrings[0]]
	if !ok {
		return rowKindID
	}
	return kind
}

func buildBaseQuery(r *http.Request, kind rowKind) query.Clauser {
	q := query.NewQuery()
	q.AddClause(selectClause(kind))
	for _, clause := range joinClauses(kind) {
		q.AddClause(clause)
	}
	q.AddClause(whereClauseBase(r))
	q.AddClause(orderClause(r))
	q.AddClause(shared.LimitClause(r))
	return q
}

func selectClause(kind rowKind) query.Clauser {
	return query.NewSelectClause("incident", rowNames(kind))
}

func rowNames(kind rowKind) []string {
	switch kind {
	case rowKindID:
		return rowNamesID[:]
	case rowKindPosition:
		return rowNamesPosition[:]
	case rowKindDetail:
		return rowNamesDetail[:]
	}
	return rowNamesID[:]
}

func joinClauses(kind rowKind) []query.Clauser {
	// TODO: Need to join on city if making a query for state IDs
	switch kind {
	case rowKindID:
	case rowKindPosition:
		return []query.Clauser{}
	case rowKindDetail:
		clauses := make([]query.Clauser, 0, len(enumTables))
		for _, table := range enumTables {
			clause := query.NewJoinClause(table)
			clauses = append(clauses, clause)
		}
		return clauses
	}
	return []query.Clauser{}
}

func whereClauseBase(r *http.Request) query.Clauser {
	w := query.NewWhereClause(query.CombinatorAnd)
	for _, table := range idQueryTables {
		clause := shared.InClause(r, table)
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

func translateFunc(kind rowKind) shared.RowTranslatorFunc {
	switch kind {
	case rowKindID:
		return translateIDRow
	case rowKindPosition:
		return translatePositionRow
	case rowKindDetail:
		return translateDetailRow
	}
	return translateIDRow
}

func translateIDRow(rows *sql.Rows) (interface{}, error) {
	row := idRow{}
	err := rows.Scan(
		&row.ID,
	)
	return row, err
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
				idRow{
					*maybe.ID,
				},
				*maybe.Name,
			}
		}
	}

	return row, err
}
