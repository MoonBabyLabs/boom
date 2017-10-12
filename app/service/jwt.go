package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/revel/modules/csrf/app"
	"io/ioutil"
	"fmt"
)

type JWT struct {
}

type KekClaim struct {
	Secret []byte
	jwt.StandardClaims
}

func (jt JWT) New() (*jwt.Token, error) {
	return jwt.New(jwt.SigningMethodHS256), nil
}

func (jt JWT) GenerateSecret() []byte {
	rnd, _ := csrf.RandomString(90)
	ioutil.WriteFile("s", []byte(rnd), 0755)

	return []byte(rnd)
}

func (jt JWT) GetSecret() ([]byte, error) {
	return ioutil.ReadFile("s")
}

func (jt JWT) Validate(tokenString string) bool {
	// Parse the token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		sec, _ := jt.GetSecret()
		return sec, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	} else {
		return false
	}
}