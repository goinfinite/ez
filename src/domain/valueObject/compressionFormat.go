package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type CompressionFormat string

var ValidCompressionFormats = []string{
	"tar", "gzip", "zip", "xz", "br",
}

func NewCompressionFormat(value interface{}) (CompressionFormat, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("CompressionFormatMustBeString")
	}

	stringValue = strings.TrimPrefix(stringValue, ".")
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidCompressionFormats, stringValue) {
		switch stringValue {
		case "gz":
			stringValue = "gzip"
		case "tarball":
			stringValue = "tar"
		case "brotli":
			stringValue = "br"
		default:
			return "", errors.New("UnsupportedCompressionFormat")
		}
	}

	return CompressionFormat(stringValue), nil
}

func (vo CompressionFormat) String() string {
	return string(vo)
}
