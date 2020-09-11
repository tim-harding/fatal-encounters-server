package shared

import (
	"fmt"
	"strings"
)

// Ordering enumerates the kinds of ordering for ORDER BY clauses
type Ordering int

const (
	// OrderingAscending sorts results in ascending order
	OrderingAscending Ordering = iota
	// OrderingDescending sorts results in descending order
	OrderingDescending
)

type orderClause struct {
	ordering Ordering
	columns  []string
}

// NewOrderClause creates a SQL ORDER BY clause
func NewOrderClause(ordering Ordering, columns []string) Clauser {
	return &orderClause{ordering, columns}
}

func (c *orderClause) String() string {
	columns := strings.Join(c.columns, ", ")
	ordering := []string{"ASC", "DESC"}[c.ordering]
	return fmt.Sprintf("ORDER BY %s %s", columns, ordering)
}

func (c *orderClause) Parameters() []interface{} {
	return []interface{}{}
}
