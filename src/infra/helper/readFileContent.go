package infraHelper

import (
	"errors"
	"os"
)

func ReadFileContent(filePath string) (string, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return "", errors.New("ReadFileInfoError: " + err.Error())
	}

	fileContentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.New("ReadFileContentError: " + err.Error())
	}

	return string(fileContentBytes), nil
}
