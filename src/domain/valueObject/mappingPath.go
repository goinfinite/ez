package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

const mappingPathRegex string = `^[^\s<>;'":#{}?\[\]]{1,512}$`

type MappingPath string

func NewMappingPath(value interface{}) (MappingPath, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("MappingPathMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)

	re := regexp.MustCompile(mappingPathRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidMappingPath")
	}

	return MappingPath(stringValue), nil
}

func (vo MappingPath) String() string {
	return string(vo)
}
