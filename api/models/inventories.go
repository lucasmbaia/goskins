package models

import (
	"fmt"
	"github.com/lucasmbaia/goskins/api/repository/filter"
)

type InventoriesFields struct {
	ID	    string  `json:",omitempty" param:"inventory"`
	User	    string  `json:",omitempty" param:"user"`
}

type Inventories struct {
	Resources
}

func (InventoriesFields) TableName() string {
	return "inventories"
}

func NewInventories() *Inventories {
	var invent = &Inventories{}

	return invent
}

func (i *Inventories) Get(filters []filter.Filters, args ...interface{}) (inventories []InventoriesFields, err error) {
	//err = i.DB.Read(filters, &inventories, args)
	fmt.Println(filters)
	return
}

func (i *Inventories) Post(data *InventoriesFields) (async bool, err error) {
	return
}
