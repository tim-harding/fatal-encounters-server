package shared

import (
	"fmt"
	"strings"
)

type inClause struct {
	column string
	values []int
}

// NewInClause creates a SQL IN clause
func NewInClause(column string, values []int) Clauser {
	return &inClause{column, values}
}

func (c *inClause) String() string {
	if len(c.values) < 1 {
		return ""
	}
	values := strings.Repeat("?, ", len(c.values))
	values = values[:len(values)-2]
	return fmt.Sprintf("%s IN (%s)", c.column, values)
}

func (c *inClause) Parameters() []interface{} {
	out := make([]interface{}, 0, len(c.values))
	for _, value := range c.values {
		out = append(out, value)
	}
	return out
}
