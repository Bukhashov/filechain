package user

import (
	"time"
	"github.com/golang-jwt/jwt/v4"
)

func (u *user) GeneratToken()(tkn string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(2*time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   "block",
		ID:        u.Dto.ID,
		Audience:  []string{"somebody_else"},
	})
	tkn, err = token.SignedString([]byte("AllYourBase"))

	return tkn, err
}