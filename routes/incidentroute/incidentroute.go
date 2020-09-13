package incidentroute

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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

// Todo: response sorting
// Todo: id queries

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
	query, err := buildIDQuery(r, kind)
	if err != nil {
		shared.Error(w, err, http.StatusBadRequest)
		return
	}
	translateRow := translateFunc(kind)
	shared.HandleRoute(w, r, query, translateRow)
}

var querystringToRowKinds = map[string]rowKind{
	"id":      rowKindID,
	"mapping": rowKindMapping,
	"listing": rowKindListing,
	"detail":  rowKindDetail,
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
	q.AddClause(whereClauseBase(r))
	q.AddClause(orderClause(r))
	q.AddClause(shared.LimitClause(r))
	return q
}

func buildIDQuery(r *http.Request, kind rowKind) (query.Clauser, error) {
	w, err := whereClauseID(r)
	if err != nil {
		return nil, err
	}
	q := query.NewQuery()
	q.AddClause(selectClause(kind))
	q.AddClause(w)
	q.AddClause(shared.LimitClause(r))
	return q, nil
}

func whereClauseID(r *http.Request) (query.Clauser, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	match := query.NewCompareClause(query.ComparisonEqual, "id", id)
	w := query.NewWhereClause(query.CombinatorAnd)
	w.AddClause(match)
	return w, nil
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

func whereClauseBase(r *http.Request) query.Clauser {
	w := query.NewWhereClause(query.CombinatorAnd)
	for _, table := range enumTables {
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

var genders = map[string]bool{
	"male":   true,
	"female": false,
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

type orderKind int

const (
	orderKindID orderKind = iota
	orderKindAge
	orderKindName
	orderKindDate
)

var orderKindColumns = []string{
	"id",
	"age",
	"name",
	"date",
}

var querystringToOrderKind = map[string]orderKind{
	"id":   orderKindID,
	"age":  orderKindAge,
	"name": orderKindName,
	"date": orderKindDate,
}

func orderClause(r *http.Request) query.Clauser {
	kind := pickOrderKind(r)
	column := orderKindColumns[kind]
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

var querystringToOrderDirection = map[string]query.Ordering{
	"ascending":  query.OrderingAscending,
	"descending": query.OrderingDescending,
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
