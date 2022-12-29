package utils

import (
	"os"
)

type FileManager interface {
	CopyNewPath(from, to string)(err error)
}

type fileManager struct {

}

func (file *fileManager) CopyNewPath(oldPath, newPath string)(err error) {
	err = os.Rename(oldPath, newPath); if err != nil {
		return err
	}
	// err = os.Remove(oldPath); if err != nil {
	// 	return err
	// }
	return nil
}

func NewFileManager() FileManager {
	return &fileManager{}
}