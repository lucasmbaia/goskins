package gorm


import (
	"time"
	"fmt"

	"github.com/lucasmbaia/goskins/api/repository/filter"
	_ "github.com/go-sql-driver/mysql"
	_gorm "github.com/jinzhu/gorm"
)

type Gorm struct {
	DB  *_gorm.DB
}

type GormConfig struct {
	Username          string  `json:",omitempty"`
	Password          string  `json:",omitempty"`
	Host              string  `json:",omitempty"`
	Port              string  `json:",omitempty"`
	DBName            string  `json:",omitempty"`
	Timeout           string  `json:",omitempty"`
	Debug             bool	  `json:",omitempty"`
	ConnsMaxIdle      int	  `json:",omitempty"`
	ConnsMaxOpen      int	  `json:",omitempty"`
	ConnsMaxLifetime  int	  `json:",omitempty"`
}

func NewGorm(cfg GormConfig) (g Gorm, err error) {
	if g.DB, err = _gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Timeout)); err != nil {
		return
	}

	g.DB.LogMode(cfg.Debug)
	g.DB.DB().SetMaxIdleConns(cfg.ConnsMaxIdle)
	g.DB.DB().SetMaxOpenConns(cfg.ConnsMaxOpen)
	g.DB.DB().SetConnMaxLifetime(time.Duration(cfg.ConnsMaxLifetime))

	return
}

func (g *Gorm) Create(entity interface{}) (err error) {
	return
}

func (g *Gorm) Read(f []filter.Filters, entity interface{}, args ...interface{}) (err error) {
	var (
		limit	  int
		operation *gorm.DB
		fields	  string
		values	  []interface{}
	)

	for i := 0; i < len(args); i++ {
		switch arg := args[i]; arg.(type) {
		case int:
			limit = arg.(int)
		}
	}

	if len(filters) > 0 {
		fields, values = filter.Join(f)
		if limit == 1 {
			operation = g.DB.Where(fields, values...).First(&entity)
		} else {
			operation = g.DB.Where(fields, values...).Find(&entity)
		}
	} else {
		operation = g.DB.Find(&entity)
	}

	if operation.Error != nil {
		err = operation.Error
	}

	return
}

func (g *Gorm) Delete(condition interface{}) (exists bool, err error) {
	return
}
