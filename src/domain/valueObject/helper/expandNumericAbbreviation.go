package voHelper

import (
	"errors"
	"strconv"
)

func ExpandNumericAbbreviation(numericStr string) (int64, error) {
	numericRegex := `(?P<number>[\d\.]+)(?P<quantifier>[KkMmBb])?`

	numericParts := FindNamedGroupsMatches(numericRegex, numericStr)
	if len(numericParts) == 0 {
		return 0, errors.New("InvalidNumericString")
	}

	numberPart := numericParts["number"]
	if numberPart == "" {
		return 0, errors.New("UndefinedNumberPart")
	}

	quantifierPart := numericParts["quantifier"]

	multiplier := float64(1)
	switch quantifierPart {
	case "K", "k":
		multiplier = 1000
	case "M", "m":
		multiplier = 1000000
	case "B", "b":
		multiplier = 1000000000
	default:
		multiplier = 1
	}

	numberPartFloat, err := strconv.ParseFloat(numberPart, 64)
	if err != nil {
		return 0, errors.New("ParseNumberError")
	}

	return int64(numberPartFloat * multiplier), nil
}
