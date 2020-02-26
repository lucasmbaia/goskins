package config

import (
	"log"
	"github.com/lucasmbaia/goskins/api/repository/gorm"
)

func LoadDB() {
	var err error

	if EnvSingletons.DB, err = gorm.NewGorm(EnvConfig.DBFields); err != nil {
		log.Fatal(err)
	}

	return
}
