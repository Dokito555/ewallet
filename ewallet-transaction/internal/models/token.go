package models

import "github.com/golang-jwt/jwt/v4"

type TokenData struct {
	UserID		int64
	Token		string
	Username	string
	FullName	string
	Email		string
}

type ClaimToken struct {	
	UserID		int		`json:"user_id"`
	Username	string	`json:"username"`
	FullName	string	`json:"fullname"`
	Email		string	`json:"email"`
	jwt.RegisteredClaims
}