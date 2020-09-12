package query

type pageClause struct {
	limit  int
	offset int
}

// NewPageClause creates a pagination clause.
// OFFSET is omitted if zero.
// Clause is omitted if limit is less than one.
func NewPageClause(limit, offset int) Clauser {
	return &pageClause{limit, offset}
}

// String returns a SQL snippet
func (p *pageClause) String() string {
	if p.limit > 0 {
		if p.offset > 0 {
			return "LIMIT ? OFFSET ?"
		}
		return "LIMIT ?"
	}
	return ""
}

// Parameters returns the SQL query placeholder contents
func (p *pageClause) Parameters() []interface{} {
	if p.limit > 0 {
		if p.offset > 0 {
			return []interface{}{
				p.limit,
				p.offset,
			}
		}
		return []interface{}{
			p.limit,
		}
	}
	return []interface{}{}
}
