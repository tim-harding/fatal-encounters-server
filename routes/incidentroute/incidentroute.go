package incidentroute

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/tim-harding/fatal-encounters-server/query"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

type coordinate struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"long"`
}

type idRow struct {
	ID int `json:"id"`
}

type mappingRow struct {
	idRow
	Coordinate coordinate `json:"coordinate"`
}

type listingRow struct {
	idRow
	Name     *string   `json:"name"`
	Age      *int      `json:"age"`
	Date     time.Time `json:"date"`
	ImageURL *string   `json:"imageUrl"`
}

type detailRow struct {
	idRow
	IsMale      *bool   `json:"isMale"`
	Zipcode     *int    `json:"zipcode"`
	Race        *int    `json:"race"`
	County      *int    `json:"county"`
	Agency      *int    `json:"agency"`
	Cause       int     `json:"cause"`
	UseOfForce  int     `json:"useOfForce"`
	City        *int    `json:"city"`
	Address     *string `json:"address"`
	Description string  `json:"description"`
	ArticleURL  *string `json:"articleUrl"`
	VideoURL    *string `json:"videoUrl"`
}

type rowKind int

const (
	rowKindID rowKind = iota
	rowKindMapping
	rowKindListing
	rowKindDetail
)

type responseRow interface {
	FromRow(rows *sql.Rows)
}

var rowNamesID = []string{
	"id",
}

var rowNamesMapping = []string{
	"id",
	"latitude",
	"longitude",
}

var rowNamesListing = []string{
	"id",
	"name",
	"age",
	"date",
	"image_url",
}

var rowNamesDetail = []string{
	"id",
	"is_male",
	"zipcode",
	"race",
	"county",
	"agency",
	"cause",
	"use_of_force",
	"city",
	"address",
	"description",
	"article_url",
	"video_url",
}

var enumTables = []string{
	"id",
	"agency",
	"cause",
	"city",
	"county",
	"race",
	"state",
	"use_of_force",
}

// HandleRouteMapping delivers basic incident information
func HandleRouteMapping(w http.ResponseWriter, r *http.Request) {
	kind := pickRowKind(r)
	buildQuery := buildQueryFactory(kind)
	translateRow := translateFunc(kind)
	shared.HandleRoute(w, r, buildQuery, translateRow)
}

func pickRowKind(r *http.Request) rowKind {
	querystrings, ok := r.URL.Query()["rowKind"]
	if !ok || len(querystrings) < 1 {
		return rowKindID
	}
	switch querystrings[0] {
	case "id":
		return rowKindID
	case "mapping":
		return rowKindMapping
	case "listing":
		return rowKindListing
	case "detail":
		return rowKindDetail
	}
	return rowKindID
}

func buildQueryFactory(kind rowKind) shared.QueryBuilderFunc {
	return func(r *http.Request) query.Clauser {
		q := query.NewQuery()
		q.AddClause(selectClause(kind))
		q.AddClause(whereClause(r))
		q.AddClause(shared.LimitClause(r))
		return q
	}
}

func selectClause(kind rowKind) query.Clauser {
	return query.NewSelectClause("incident", rowNames(kind))
}

func rowNames(kind rowKind) []string {
	switch kind {
	case rowKindID:
		return rowNamesID
	case rowKindMapping:
		return rowNamesMapping
	case rowKindListing:
		return rowNamesListing
	case rowKindDetail:
		return rowNamesDetail
	}
	return rowNamesID
}

func whereClause(r *http.Request) query.Clauser {
	w := query.NewWhereClause(query.CombinatorAnd)
	for _, table := range enumTables {
		clause := shared.InClause(r, table)
		w.AddClause(clause)
	}
	w.AddClause(shared.SearchClause(r))
	w.AddClause(ageClause(r, "ageMin", query.ComparatorGreaterEqual))
	w.AddClause(ageClause(r, "ageMax", query.ComparatorLesserEqual))
	w.AddClause(genderMaskClause(r))
	w.AddClause(dateMaskClause(r, "dateMin", query.ComparatorGreaterEqual))
	w.AddClause(dateMaskClause(r, "dateMax", query.ComparatorLesserEqual))
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
	if !ok || len(querystrings) < 1 {
		return nil
	}
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

func dateMaskClause(r *http.Request, key string, comparator query.Comparator) query.Clauser {
	querystrings, ok := r.URL.Query()[key]
	if !ok || len(querystrings) < 1 {
		return nil
	}
	t, err := time.Parse("2013-Feb-03", querystrings[0])
	if err != nil {
		return nil
	}
	return query.NewCompareClause(comparator, "date", t)
}

func translateFunc(kind rowKind) shared.RowTranslatorFunc {
	switch kind {
	case rowKindID:
		return translateIDRow
	case rowKindMapping:
		return translateMappingRow
	case rowKindListing:
		return translateListingRow
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

func translateMappingRow(rows *sql.Rows) (interface{}, error) {
	row := mappingRow{}
	err := rows.Scan(
		&row.ID,
		&row.Coordinate.Latitude,
		&row.Coordinate.Longitude,
	)
	return row, err
}

func translateListingRow(rows *sql.Rows) (interface{}, error) {
	row := listingRow{}
	err := rows.Scan(
		&row.ID,
		&row.Name,
		&row.Age,
		&row.Date,
		&row.ImageURL,
	)
	return row, err
}

func translateDetailRow(rows *sql.Rows) (interface{}, error) {
	row := detailRow{}
	err := rows.Scan(
		&row.ID,
		&row.IsMale,
		&row.Zipcode,
		&row.Race,
		&row.County,
		&row.Agency,
		&row.Cause,
		&row.UseOfForce,
		&row.City,
		&row.Address,
		&row.Description,
		&row.ArticleURL,
		&row.VideoURL,
	)
	return row, err
}
