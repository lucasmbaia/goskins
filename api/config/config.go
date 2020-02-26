package config

import (
	"github.com/lucasmbaia/goskins/api/repository/broker"
	"github.com/lucasmbaia/goskins/api/repository/gorm"
	"github.com/lucasmbaia/goskins/steam-api"
)

var (
	EnvSingletons	Singletons
	EnvConfig	Config
)

type Singletons struct {
	DB	broker.Brokers
	Session	map[string]steam.Session
}

type Config struct {
	DBFields    gorm.GormConfig `json:",omitempty"`
}
