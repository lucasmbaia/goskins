package models

import (
	"github.com/lucasmbaia/goskins/api/repository/filter"
)

type UsersFields struct {
	ID	    string  `json:",omitempty" url:"User"`
	Name	    string  `json:",omitempty"`
	NickName    string  `json:",omitempty"`
	SteamID	    string  `json:"-"`
}

type Users struct {
	Resources
}

func NewUsers() *Users {
	var users = &Users{}

	return users
}

func (u *Users) Get(filters []filter.Filters, args ...interface{}) (users []UsersFields, err error) {
	//err = u.DB.Read(filters, &users, args)
	users = []UsersFields{{
		ID:	    "e963f7a2-cdce-4514-9204-9a670811c704",
		Name:	    "teste",
		NickName:   "teste",
		SteamID:    "1",
	}}

	return
}

func (u *Users) Post(data *UsersFields) (async bool, err error) {
	data.ID = "e963f7a2-cdce-4514-9204-9a670811c704"
	return
}
