package filter

import (
	"fmt"
)

const (
	EQ    = "="
	NEQ   = "<>"
	LIKE  = "LIKE"
	IN    = "IN"
)

type Filters struct {
	Conditions  Conditions
}

type Conditions struct {
	Field string
	Value interface{}
	op    string
}

func (c Conditions) Eq(v interface{}) Conditions {
	c.op = EQ
	c.Value = v

	return c
}

func (c Conditions) Neq(v interface{}) Conditions {
	c.op = NEQ
	c.Value = v

	return c
}
func (c Conditions) Like(v interface{}) Conditions {
	c.op = LIKE
	c.Value = v

	return c
}
func (c Conditions) In(v interface{}) Conditions {
	c.op = IN
	c.Value = v

	return c
}

func Join(f []Filters) (c string, v []interface{}) {
	for _, filter := range f {
		switch filter.Conditions.op {
		case EQ:
			if c == "" {
				c += fmt.Sprintf("%s = ?", filter.Conditions.Field)
			} else {
				c += fmt.Sprintf("AND %s = ?", filter.Conditions.Field)
			}
		case NEQ:
			if c == "" {
				c += fmt.Sprintf("%s <> ?", filter.Conditions.Field)
			} else {
				c += fmt.Sprintf("AND %s <> ?", filter.Conditions.Field)
			}
		case LIKE:
			if c == "" {
				c += fmt.Sprintf("%s LIKE ?", filter.Conditions.Field)
			} else {
				c += fmt.Sprintf("AND %s LIKE ?", filter.Conditions.Field)
			}
		case IN:
			if c == "" {
				c += fmt.Sprintf("%s IN (?)", filter.Conditions.Field)
			} else {
				c += fmt.Sprintf("AND %s IN (?)", filter.Conditions.Field)
			}
		}

		v = append(v, filter.Conditions.Value)
	}

	return
}
