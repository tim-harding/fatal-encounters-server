package query

import "fmt"

type conditionsClause struct {
	expr Subclauser
}

// Combinator is an enumeration of SQL WHERE combinators
type Combinator int

const (
	// CombinatorAnd is a SQL AND keyword
	CombinatorAnd Combinator = iota
	// CombinatorOr is a SQL OR keyword
	CombinatorOr
)

var tokens = []string{
	" AND ",
	" OR ",
}

// NewConditionsClause creates a new grouping of clauses for a WHERE statement
func NewConditionsClause(combinator Combinator) Subclauser {
	combinatorToken := tokens[combinator]
	expr := newSubexpression(combinatorToken)
	return &conditionsClause{expr}
}

func (c *conditionsClause) String() string {
	inner := c.expr.String()
	if len(inner) < 1 {
		return ""
	}
	return fmt.Sprintf("(%s)", inner)
}

func (c *conditionsClause) Parameters() []interface{} {
	return c.expr.Parameters()
}

func (c *conditionsClause) AddClause(clause Clauser) {
	c.expr.AddClause(clause)
}
