package infraHelper

import "encoding/base64"

func EncodeStr(decodedStr string) string {
	return base64.StdEncoding.EncodeToString([]byte(decodedStr))
}

func DecodeStr(encodedStr string) (string, error) {
	encodedStrBytes, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		return "", err
	}

	return string(encodedStrBytes), nil
}
