package common

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("te")

type Claims struct {
	ID       uint
	Email    string
	Password string
	Name     string
	Role     string
	Phone    string
	jwt.StandardClaims
}

func ReleaseToken(claims Claims) (string, error) {
	//TODO:发放token不需要那么长时间，按照需求来修改
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "TODO",
		Subject:   "TODO",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}
