package query

import "fmt"

type joinClause struct {
	table string
}

// NewJoinClause creates a new join clause
func NewJoinClause(table string) Clauser {
	return &joinClause{table}
}

func (j *joinClause) String() string {
	return fmt.Sprintf("JOIN %s ON %s_id=%s.id", j.table, j.table, j.table)
}

func (j *joinClause) Parameters() []interface{} {
	return []interface{}{}
}
