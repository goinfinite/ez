package uiHelper

import (
	"fmt"
)

func FormatPointer[ParamType interface{}](
	pointer *ParamType, nilPlaceholder string,
) string {
	if pointer == nil {
		return nilPlaceholder
	}

	return fmt.Sprintf("%v", *pointer)
}
