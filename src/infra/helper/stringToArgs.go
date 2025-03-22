package infraHelper

import (
	"regexp"
	"strings"
)

func StringToArgs(inputString string) []string {
	argsRegex := regexp.MustCompile(`"[^"]*"|'[^']*'|[^ ]+`)
	rawArgs := argsRegex.FindAllString(inputString, -1)
	for argIndex, rawArg := range rawArgs {
		rawArgWithoutQuoting := rawArg

		rawArgWithoutQuoting = strings.TrimPrefix(rawArgWithoutQuoting, "'")
		rawArgWithoutQuoting = strings.TrimSuffix(rawArgWithoutQuoting, "'")
		rawArgWithoutQuoting = strings.TrimPrefix(rawArgWithoutQuoting, "\"")
		rawArgWithoutQuoting = strings.TrimSuffix(rawArgWithoutQuoting, "\"")

		rawArgs[argIndex] = rawArgWithoutQuoting
	}
	return rawArgs
}
