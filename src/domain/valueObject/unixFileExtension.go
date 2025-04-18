package valueObject

import (
	"errors"
	"mime"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const unixFileExtensionRegexExpression = `^([\w\-]{1,15}\.)?[\w\-]{1,15}$`

type UnixFileExtension string

func NewUnixFileExtension(value interface{}) (
	unixFileExtension UnixFileExtension, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return unixFileExtension, errors.New("UnixFileExtensionMustBeString")
	}
	stringValue = strings.TrimPrefix(stringValue, ".")

	re := regexp.MustCompile(unixFileExtensionRegexExpression)
	if !re.MatchString(stringValue) {
		return unixFileExtension, errors.New("InvalidUnixFileExtension")
	}

	return UnixFileExtension(stringValue), nil
}

func (vo UnixFileExtension) ReadMimeType() MimeType {
	mimeTypeStr := "generic"

	fileExtWithLeadingDot := "." + string(vo)
	mimeTypeWithCharset := mime.TypeByExtension(fileExtWithLeadingDot)
	if len(mimeTypeWithCharset) > 1 {
		mimeTypeOnly := strings.Split(mimeTypeWithCharset, ";")[0]
		mimeTypeStr = mimeTypeOnly
	}

	mimeType, _ := NewMimeType(mimeTypeStr)
	return mimeType
}

func (vo UnixFileExtension) String() string {
	return string(vo)
}
