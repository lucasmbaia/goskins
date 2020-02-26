package models

import (
	"fmt"
	"github.com/lucasmbaia/goskins/api/repository/filter"
)

type InventoriesFields struct {
	ID	    string  `json:",omitempty"`
	User	    string  `json:",omitempty"`
}

type Inventories struct {
	Resources
}

func NewInventories() *Inventories {
	var invent = &Inventories{}

	return invent
}

func (i *Inventories) Get(filters []filter.Filters, args ...interface{}) (inventories []InventoriesFields, err error) {
	err = i.DB.Read(filters, &inventories, args)
	return
}

func (i *Inventories) Post(data *InventoriesFields) (async bool, err error) {
	fmt.Println(data)
	return
}
