package common


import (
	"server/model"
	"github.com/dgrijalva/jwt-go"

	"time"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	//UserId uint
	Uuid string
	Password string
	jwt.StandardClaims
}

func ReleaseUserToken(user model.User) (string, error)  {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Uuid:        user.Uuid,
		Password:        user.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "clf.model",
			Subject: "user.token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseUserToken(tokenString string) (*jwt.Token, *Claims, error){
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}

func ReleaseAdminToken(admin model.Admin) (string, error)  {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		Uuid:        admin.Uuid,
		Password:        admin.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "clf.model",
			Subject: "user.token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseAdminToken(tokenString string) (*jwt.Token, *Claims, error){
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}
