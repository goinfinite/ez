package infraHelper

import (
	"errors"
	"log"
)

func DownloadFile(url string, filePath string) error {
	_, err := RunCmd(
		"wget",
		"-q",
		"--no-check-certificate",
		"-O",
		filePath,
		url,
	)

	if err != nil {
		log.Printf("DownloadFileError: %s", err)
		return errors.New("DownloadFileError")
	}

	return nil
}
