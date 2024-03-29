package dto

import (
	"mime/multipart"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	jwt.RegisteredClaims
	ID 		int64 `json:"id"`
	Name 	string `json:"name"`
	Email	string `json:"email"`
	Image	string `json:"image"`
	File	*multipart.FileHeader `json:"file"`
}