package controllers

import (
	"github.com/lucasmbaia/goskins/api/models/interfaces"
	"github.com/lucasmbaia/goskins/api/models"
)

type Inventories struct {
	Resources
}

func NewInventories() *Inventories {
	return &Inventories{
		Resources{
			GetModel: func() interfaces.Models {
				return models.NewResources(models.NewInventories())
			},
			GetFields: func() interface{} {
				return &models.InventoriesFields{}
			},
		},
	}
}
