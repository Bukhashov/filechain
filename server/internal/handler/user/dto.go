package user

import (
	"mime/multipart"
)

type Dto struct {
	ID 		string `json:"id"`
	Name 	string `json:"name"`
	Email	string `json:"email"`
	Image	string `json:"image"`
	File	*multipart.FileHeader `json:"file"`
}