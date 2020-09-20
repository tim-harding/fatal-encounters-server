package query

import "fmt"

type insertClause struct {
	table string
}

// NewInsertClause creates a new INSERT INTO clause
func NewInsertClause(table string) Clauser {
	return &insertClause{table}
}

func (i *insertClause) String() string {
	return fmt.Sprintf("INSERT INTO %s", i.table)
}

func (i *insertClause) Parameters() []interface{} {
	return []interface{}{}
}
