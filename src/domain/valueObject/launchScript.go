package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type LaunchScript string

func NewLaunchScript(value interface{}) (LaunchScript, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("LaunchScriptMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)

	if len(stringValue) < 3 {
		return "", errors.New("LaunchScriptIsTooShort")
	}

	if len(stringValue) > 16384 {
		return "", errors.New("LaunchScriptIsTooLong")
	}

	hasSheBang := strings.HasPrefix(stringValue, "#!")
	if !hasSheBang {
		stringValue = "#!/bin/sh\n" + stringValue
	}

	return LaunchScript(stringValue), nil
}

func NewLaunchScriptFromEncodedContent(
	encodedContent EncodedContent,
) (LaunchScript, error) {
	var launchScript LaunchScript

	decodedContent, err := encodedContent.GetDecoded()
	if err != nil {
		return launchScript, err
	}

	return NewLaunchScript(decodedContent)
}

func (ls LaunchScript) String() string {
	return string(ls)
}
