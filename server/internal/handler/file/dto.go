package file

import (
	"mime/multipart"
)

type Dto struct {
	Username 	string 					`json:"username"`
	Image		string 					`json:"image"`
	Title		string					`json:"title"`
	Address		string					`json:"address"`
	File		*multipart.FileHeader 	`json:"file"`
	FilePath	string					`json:"filepath"`
}