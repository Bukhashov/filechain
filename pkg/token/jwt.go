package token

import (
	// "context"
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"github.com/Bukhashov/filechain/configs"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/dto"
	// "github.com/jackc/pgx/v5/pgxpool"
	// "github.com/Bukhashov/filechain/pkg/logging"
)

type Token interface {
	Generator(u *model.User)(newToken string, err error)
	Parse(clientToken string, d *dto.User)(err error)
	// Control(ctx context.Context, client *pgxpool.Pool, logger *logging.Logger)(err error)
}

type token struct {
	key 	string
}

func (t *token) Generator(u *model.User)(newToken string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, dto.User{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2*time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "filechain",
		},
		ID: u.ID,
		Name: u.Name,
		Email: u.Email,
	})

	newToken, err = jwtToken.SignedString([]byte(t.key)); if err != nil {
		return "", err
	}
	return newToken, nil
}

func (t *token) Parse(clientToken string, d *dto.User)(err error){
	cfg := configs.GetConfig()
	token, err :=  jwt.ParseWithClaims(clientToken, d, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Token.Key), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}
// func (t *token) Control(ctx context.Context, client *pgxpool.Pool, logger *logging.Logger)(err error){
// 	s := NewStorage(client, logger)
	
// 	err = s.FindUserByEmail(ctx, &model.User{Email: d.Email}); if err != nil {
// 		return errors.New("user token not fund")
// 	}
// 	return nil
// }

func NewToken(key string) Token {
	return &token{
		key: key,
	}
}