package models

import (
	"time"
)

//User defines what a user looks like
type User struct {
	Id 			uint			`json:"id"`
	Name 		string			`json:"name"`
	Password 	string			`json:"password"`
	Email 		string			`json:"email"`
	Address 	string			`json:"address"`
	Created 	time.Time		`json:"created"`
	Roles		[]Role			`json:"roles"`
	Permissions	[]Permission	`json:"permissions"`
}