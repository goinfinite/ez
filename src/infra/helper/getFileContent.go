package infraHelper

import (
	"errors"
	"os"
)

func GetFileContent(filePath string) (string, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return "", errors.New("GetFileInfoError: " + err.Error())
	}

	fileContentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.New("GetFileContentError: " + err.Error())
	}

	return string(fileContentBytes), nil
}
