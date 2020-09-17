package incidentroute

import "github.com/tim-harding/fatal-encounters-server/query"

func selectClause(kind rowKind) query.Clauser {
	return query.NewSelectClause("incident", rowNames[kind])
}
