package uiHelper

import (
	"fmt"
)

func FormatPointer[ParamType interface{}](pointer *ParamType) string {
	if pointer == nil {
		return "-"
	}

	return fmt.Sprintf("%v", *pointer)
}
