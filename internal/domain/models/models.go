package models

type User struct {
	UID 		int 	`json:"uid"`
	Name 		string 	`json:"name"`
	Login 		string 	`json:"login"`
	Password 	string 	`json:"password"`
}

type Response struct {
	Message string
}