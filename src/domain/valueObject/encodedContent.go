package valueObject

import (
	"encoding/base64"
	"errors"
	"regexp"
)

const encodedContentRegexExpression = `^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{4}|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{2}={2})$`

type EncodedContent string

func NewEncodedContent(value string) (EncodedContent, error) {
	if len(value) == 0 {
		return "", errors.New("EmptyEncodedContent")
	}

	ecRegex := regexp.MustCompile(encodedContentRegexExpression)
	isValid := ecRegex.MatchString(value)
	if !isValid {
		return "", errors.New("InvalidEncodedContent")
	}

	return EncodedContent(value), nil
}

func (ec EncodedContent) GetDecoded() (string, error) {
	decodedContent, err := base64.StdEncoding.DecodeString(string(ec))
	if err != nil {
		return "", err
	}

	return string(decodedContent), nil
}

func (ec EncodedContent) String() string {
	return string(ec)
}
