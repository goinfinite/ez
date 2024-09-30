package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type LaunchScript string

func NewLaunchScript(value interface{}) (launchScript LaunchScript, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return launchScript, errors.New("LaunchScriptMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)

	if len(stringValue) < 3 {
		return launchScript, errors.New("LaunchScriptIsTooShort")
	}

	if len(stringValue) > 16384 {
		return launchScript, errors.New("LaunchScriptIsTooLong")
	}

	hasSheBang := strings.HasPrefix(stringValue, "#!")
	if !hasSheBang {
		stringValue = "#!/bin/sh\n" + stringValue
	}

	return LaunchScript(stringValue), nil
}

func NewLaunchScriptFromEncodedContent(
	encodedContent EncodedContent,
) (launchScript LaunchScript, err error) {
	decodedContent, err := encodedContent.GetDecoded()
	if err != nil {
		return launchScript, err
	}

	return NewLaunchScript(decodedContent)
}

func (vo LaunchScript) String() string {
	return string(vo)
}
