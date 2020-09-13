package incidentroute

import "github.com/tim-harding/fatal-encounters-server/query"

// Miscellaneous
// ------------------------------------------------------------

var (
	enumTables = [...]string{
		"id",
		"agency",
		"cause",
		"city",
		"county",
		"race",
		"state",
		"use_of_force",
	}

	genders = map[string]bool{
		"male":   true,
		"female": false,
	}

	querystringToOrderDirection = map[string]query.Ordering{
		"ascending":  query.OrderingAscending,
		"descending": query.OrderingDescending,
	}
)

// Row names
// ------------------------------------------------------------

var (
	rowNamesID = [...]string{
		"id",
	}

	rowNamesMapping = [...]string{
		"id",
		"latitude",
		"longitude",
	}

	rowNamesListing = [...]string{
		"id",
		"name",
		"age",
		"date",
		"image_url",
	}

	rowNamesDetail = [...]string{
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
)

// Row kind
// ------------------------------------------------------------

type rowKind int

const (
	rowKindID rowKind = iota
	rowKindMapping
	rowKindListing
	rowKindDetail
)

var (
	querystringToRowKinds = map[string]rowKind{
		"id":      rowKindID,
		"mapping": rowKindMapping,
		"listing": rowKindListing,
		"detail":  rowKindDetail,
	}
)

// Order kind
// ------------------------------------------------------------

type orderKind int

const (
	orderKindID orderKind = iota
	orderKindAge
	orderKindName
	orderKindDate
)

var (
	orderKindColumns = [...]string{
		"id",
		"age",
		"name",
		"date",
	}

	querystringToOrderKind = map[string]orderKind{
		"id":   orderKindID,
		"age":  orderKindAge,
		"name": orderKindName,
		"date": orderKindDate,
	}
)
