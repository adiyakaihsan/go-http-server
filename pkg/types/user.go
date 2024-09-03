package types

import "github.com/golang-jwt/jwt/v4"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       int    `json:"id"`
}

type Claims struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	jwt.StandardClaims
}

type TokenResponse struct {
    Token string `json:"token"`
}