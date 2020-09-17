package query

import "fmt"

type notClause struct {
	inner Clauser
}

// NewNotClause creates a new SQL NOT clause
func NewNotClause(inner Clauser) Clauser {
	return &notClause{inner}
}

func (n *notClause) String() string {
	inner := n.inner.String()
	if len(inner) == 0 {
		return ""
	}
	return fmt.Sprintf("NOT %s", n.inner.String())
}

func (n *notClause) Parameters() []interface{} {
	return n.inner.Parameters()
}
