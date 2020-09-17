package incidentroute

import "github.com/tim-harding/fatal-encounters-server/query"

// Miscellaneous
// ------------------------------------------------------------

var (
	enumTables = [...]string{
		"agency",
		"cause",
		"city",
		"county",
		"race",
		"use_of_force",
	}

	idQueryTables = [...]string{
		// Same as enumTables
		"agency",
		"cause",
		"city",
		"county",
		"race",
		"use_of_force",
		// ...plus state
		"state",
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
	rowNames = [...][]string{
		{
			"incident.id",
		},
		{
			"incident.id",
			"incident.latitude",
			"incident.longitude",
		},
		{
			"incident.id",
			"incident.name",
			"incident.age",
			"incident.date",
			"incident.image_url",
			"incident.is_male",
			"incident.address",
			"incident.description",
			"incident.article_url",
			"incident.video_url",
			"incident.zipcode",

			"cause.id",
			"cause.name",

			"use_of_force.id",
			"use_of_force.name",

			"race.id",
			"race.name",

			"county.id",
			"county.name",

			"agency.id",
			"agency.name",

			"city.id",
			"city.name",
		},
	}
)

type rowKind int

const (
	rowKindFilter rowKind = iota
	rowKindPosition
	rowKindDetail
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
