package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

func genToken(userId int) (string, error) {
	var secretkey = []byte("f458a57132a1c50ce5064937a10ed33bf27086ae89598daf33fb620d99bdb95b")

	expireDate := time.Now().Add(100 * time.Minute)
	claims := &Claims{
		Id: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireDate.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstring, err := token.SignedString(secretkey)
	if err != nil {
		return "", err
	}
	return tokenstring, nil

}
