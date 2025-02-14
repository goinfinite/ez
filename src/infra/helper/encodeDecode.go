package infraHelper

import "encoding/base64"

func EncodeStr(decodedStr string) (string, error) {
	encodedStr, err := base64.StdEncoding.DecodeString(decodedStr)
	if err != nil {
		return "", err
	}

	return string(encodedStr), nil
}

func DecodeStr(encodedStr string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(encodedStr)), nil
}
