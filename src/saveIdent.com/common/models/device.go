package models

type Device struct {
	Id         uint `json:"id"`
	User       User `json:"user"`
	IsReceiver bool `json:"is_receiver"`
}