package query

import "strings"

type subexpression struct {
	connector string
	parts     []Clauser
}

// NewSubexpression joins together SQL clauses
func newSubexpression(connector string) subexpression {
	return subexpression{connector, []Clauser{}}
}

func (s *subexpression) String() string {
	parts := make([]string, 0, len(s.parts))
	for _, part := range s.parts {
		text := part.String()
		if len(text) > 0 {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, s.connector)
}

func (s *subexpression) Parameters() []interface{} {
	parameters := make([]interface{}, 0)
	for _, query := range s.parts {
		for _, clause := range query.Parameters() {
			parameters = append(parameters, clause)
		}
	}
	return parameters
}

func (s *subexpression) AddClause(clause Clauser) {
	if clause != nil {
		s.parts = append(s.parts, clause)
	}
}
