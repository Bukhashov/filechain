package utils

import (
	"os"
)

func CopyNewPath(oldPath, newPath string)(err error) {
	err = os.Rename(oldPath, newPath); if err != nil {
		return err
	}
	// err = os.Remove(oldPath); if err != nil {
	// 	return err
	// }
	return nil
}