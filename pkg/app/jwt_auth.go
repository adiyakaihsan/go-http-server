package app

import (
	"fmt"
	"log"
	"time"

	"github.com/adiyakaihsan/go-http-server/pkg/config"
	"github.com/adiyakaihsan/go-http-server/pkg/types"
	jwt "github.com/golang-jwt/jwt/v4"
)

func getDefaultSigningMethod(method string) jwt.SigningMethod {
	if method != "" {
		return jwt.GetSigningMethod(method)
	} else {
		return jwt.SigningMethodHS256
	}
}

func generateJWTToken(id int, username string) (types.TokenResponse, error) {

	claims := types.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    config.APP_NAME,
			ExpiresAt: time.Now().Add(config.LOGIN_EXPIRATION_DURATION).Unix(),
		},
		ID:       id,
		Username: username,
	}

	token := jwt.NewWithClaims(
		getDefaultSigningMethod(config.Jwt_signing_method),
		claims,
	)
	log.Printf("%s: %v", "Token", token)
	signedToken, err := token.SignedString([]byte(config.JWT_SIGNATURE_KEY))
	if err != nil {
		log.Printf("%s: %v", "Error when generating token", err)
		return types.TokenResponse{}, err
	}

	tokenResponse := types.TokenResponse{Token: signedToken}
	log.Printf("%s: %v", "isi tokenResponse", tokenResponse)

	return tokenResponse, nil
}

func parseAuthToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("signing method invalid")
		}

		return []byte(config.JWT_SIGNATURE_KEY), nil
	})

	if err != nil {
		log.Printf("%s: %v", "Error when parsing token", err)
		return nil, fmt.Errorf("error when parsing token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Printf("%s: %v", "Error when validating token", err)
		return nil, fmt.Errorf("error when validating token")
	}

	return claims, nil

}
