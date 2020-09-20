package query

import "fmt"

type groupClause struct {
	Column string
}

// NewGroupClause creates a new GROUP BY clause
func NewGroupClause(column string) Clauser {
	return &groupClause{column}
}

func (g *groupClause) String() string {
	return fmt.Sprintf("GROUP BY %s", g.Column)
}

func (g *groupClause) Parameters() []interface{} {
	return []interface{}{}
}
