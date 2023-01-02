package user

import (
	"errors"
	"time"
	"fmt"

	"github.com/Bukhashov/filechain/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

func (u *user) GeneratorJWT() (err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Dto{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2*time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "filechain",
		},
		ID: u.Model.ID,
		Name: u.Model.Name,
		Email: u.Model.Email,
	})

	u.Token.jwt, err = token.SignedString([]byte(u.config.Token.Key)); if err != nil {
		return err
	}

	return nil
}

func (d *Dto) ParseJwt(clientToken string)(err error){
	cfg := config.GetConfig()

	token, err :=  jwt.ParseWithClaims(clientToken, d, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Token.Key), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	fmt.Printf("", )
	return nil
}