package types

import "github.com/golang-jwt/jwt/v4"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       int    `json:"id"`
}

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	jwt.StandardClaims
}

type TokenResponse struct {
	Token string `json:"token"`
}

type contextKey string

const UserInfoKey contextKey = "userInfo"
