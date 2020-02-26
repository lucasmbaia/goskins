package decorators

import (
	"github.com/lucasmbaia/goskins/api/models/interfaces"
)

type Transaction struct {
	interfaces.Models
}

func NewTransaction(m interfaces.Models) interfaces.Models {
	return &Transaction{m}
}

func (t *Transaction) Get(data interface{}) (response interface{}, err error) {
	response, err = t.Models.Get(data)
	return
}

func (t *Transaction) Post(data interface{}) (async bool, err error) {
	async, err = t.Models.Post(data)
	return
}
