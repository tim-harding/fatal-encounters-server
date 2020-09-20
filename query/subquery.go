package query

import "fmt"

type subquery struct {
	expr Clauser
}

// NewSubquery adds parentheses around a clause
func NewSubquery(expr Clauser) Clauser {
	return &subquery{expr}
}

func (s *subquery) String() string {
	return fmt.Sprintf("(%s)", s.expr.String())
}

func (s *subquery) Parameters() []interface{} {
	return s.expr.Parameters()
}
