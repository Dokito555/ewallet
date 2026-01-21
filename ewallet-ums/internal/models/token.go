package models

import "github.com/golang-jwt/jwt/v4"

type RefreshTokenResponse struct {
	Token		string	`json:"token"`
}	

type ClaimToken struct {	
	UserID		int		`json:"user_id"`
	Username	string	`json:"username"`
	FullName	string	`json:"fullname"`
	Email		string	`json:"email"`
	jwt.RegisteredClaims
}