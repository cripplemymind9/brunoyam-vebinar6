package models

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	UID 		int 	`json:"uid"`
	Name 		string 	`json:"name" validate:"required"`
	Login 		string 	`json:"login" validate:"required"`
	Password 	string 	`json:"password" validate:"required"`
}

type LoginUser struct {
	Login 		string 	`json:"login" validate:"required"`
	Password 	string 	`json:"password" validate:"required"`
}

type Claims struct {
	UID 		int 	`json:"uid"`
	Login 		string 	`json:"login"`
	jwt.StandardClaims
}

type Response struct {
	Message string
}