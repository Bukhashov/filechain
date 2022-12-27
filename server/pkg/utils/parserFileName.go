package utils

import (
	"path/filepath"
)

func ParseFileName(file string)(name, extension string){
	extension = filepath.Ext(file)
	name = file[:len(file)-len(extension)]
	return name, extension
}