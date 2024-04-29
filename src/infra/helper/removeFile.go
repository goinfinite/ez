package infraHelper

import "os"

func RemoveFile(filePath string) error {
	_, err := os.Stat(filePath)
	fileExists := !os.IsNotExist(err)
	if !fileExists {
		return nil
	}

	return os.Remove(filePath)
}
