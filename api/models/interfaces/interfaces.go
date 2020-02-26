package interfaces

import (
	"github.com/lucasmbaia/goskins/api/repository/filter"
)

type Models interface {
	//Get(interface{}) (interface{}, error)
	Get([]filter.Filters, ...interface{}) (interface{}, error)
	Post(interface{}) (bool, error)
}
