package infraHelper

import (
	"os"
)

func MakeDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	return os.MkdirAll(dirPath, os.ModePerm)
}
