package query

type rawSQL struct {
	text string
}

// NewRawSQL creates a clause containing the given SQL query text
func NewRawSQL(text string) Clauser {
	return &rawSQL{text}
}

func (r *rawSQL) String() string {
	return r.text
}

func (r *rawSQL) Parameters() []interface{} {
	return []interface{}{}
}
