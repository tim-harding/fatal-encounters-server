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
		// Need to join on city for this
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
	rowNamesID = [...]string{
		"incident.id",
	}

	rowNamesPosition = [...]string{
		"incident.id",
		"incident.latitude",
		"incident.longitude",
	}

	rowNamesDetail = [...]string{
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
	}
)

// Row kind
// ------------------------------------------------------------

type rowKind int

const (
	rowKindID rowKind = iota
	rowKindPosition
	rowKindDetail
)

var (
	querystringToRowKinds = map[string]rowKind{
		"id":       rowKindID,
		"position": rowKindPosition,
		"detail":   rowKindDetail,
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
