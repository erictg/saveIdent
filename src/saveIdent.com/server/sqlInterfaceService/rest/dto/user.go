package dto

type RegisterRequestDTO struct{
	Name 		string		`json:"name"`
	Password 	string		`json:"password"`
	Address		string		`json:"address"`
	Email		string		`json:"email"`
}