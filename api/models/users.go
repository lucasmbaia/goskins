package models

import (
	"github.com/lucasmbaia/goskins/api/repository/filter"
	"github.com/lucasmbaia/goskins/api/config"
)

type UsersFields struct {
	ID	    string  `json:",omitempty" param:"user" model:"id"`
	Name	    string  `json:",omitempty" model:"name"`
	NickName    string  `json:",omitempty" model:"nick_name"`
	TraderURL   string  `json:",omitempty" model:"trader_url"`
	SteamID	    string  `json:"-" model:"steam_id"`
}

func (UsersFields) TableName() string {
	return "users"
}

type Users struct {
	Resources
}

func NewUsers() *Users {
	var users = &Users{}
	users.DB = config.EnvSingletons.DB

	return users
}

func (u *Users) Get(filters []filter.Filters, args ...interface{}) (users []UsersFields, err error) {
	users = []UsersFields{}
	err = u.DB.Read(filters, &users, args)

	return
}

func (u *Users) Post(data *UsersFields) (async bool, err error) {
	return
}

func (u *Users) Patch(data *UsersFields, filters []filter.Filters) (err error) {
	return
}
