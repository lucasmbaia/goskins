package models

import (
	"reflect"

	"github.com/lucasmbaia/goskins/api/models/interfaces"
	"github.com/lucasmbaia/goskins/api/models/decorators"
	"github.com/lucasmbaia/goskins/api/repository/broker"
	"github.com/lucasmbaia/goskins/api/repository/filter"
)

type Resources struct {
	Model	reflect.Value
	DB	broker.Brokers
}

func NewResources(m interface{}) interfaces.Models {
	var model = reflect.ValueOf(m)

	return decorators.NewTransaction(&Resources{Model: model})
}

func (r *Resources) Get(filters []filter.Filters, a ...interface{}) (response interface{}, err error) {
	var (
		args	[]reflect.Value
		output	[]reflect.Value
		ok	bool
	)

	args = append(args, reflect.ValueOf(filters))
	args = append(args, reflect.ValueOf(a))
	output = r.Model.MethodByName("Get").Call(args)

	if err, ok = output[1].Interface().(error); ok {
		return
	}

	response = output[0].Interface()

	return
}

/*func (r *Resources) Get(data interface{}) (response interface{}, err error) {
	var (
		args	[]reflect.Value
		output	[]reflect.Value
		v	reflect.Value
		ok	bool
		f	[]filter.Filters
		a	[]interface{}
		t	reflect.Type
	)

	v = reflect.ValueOf(data).Elem()
	t = v.Type()

	for i := 0; i < v.NumField(); i++ {
		f = append(f, filter.Filters{
			Conditions: filter.Conditions{
				Field:	t.Field(i).Name,
				Value:	v.Field(i).Interface(),
			},
		})
	}

	args = append(args, reflect.ValueOf(f))
	args = append(args, reflect.ValueOf(a))
	output = r.Model.MethodByName("Get").Call(args)

	if err, ok = output[1].Interface().(error); ok {
		return
	}

	response = output[0].Interface()

	return
}*/

func (r *Resources) Post(data interface{}) (async bool, err error) {
	var (
		args	[]reflect.Value
		output	[]reflect.Value
		v	reflect.Value
		ok	bool
	)

	v = reflect.ValueOf(data)
	args = append(args, v)
	output = r.Model.MethodByName("Post").Call(args)

	if err, ok = output[1].Interface().(error); ok {
		return
	}

	async = output[0].Interface().(bool)

	return
}
