package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type MappingMatchPattern string

var ValidMappingMatchPatterns = []string{
	"begins-with",
	"contains",
	"equals",
	"ends-with",
}

func NewMappingMatchPattern(value interface{}) (MappingMatchPattern, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("MatchPatternMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidMappingMatchPatterns, stringValue) {
		return "", errors.New("InvalidMappingMatchPattern")
	}

	return MappingMatchPattern(stringValue), nil
}

func (vo MappingMatchPattern) String() string {
	return string(vo)
}
