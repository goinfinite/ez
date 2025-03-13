package sharedHelper

import (
	"log/slog"
	"reflect"
	"strings"
)

func StringSliceValueObjectParser[TypedObject any](
	rawInputValues any,
	valueObjectConstructor func(any) (TypedObject, error),
) []TypedObject {
	parsedObjects := make([]TypedObject, 0)

	if rawInputValues == nil {
		return parsedObjects
	}

	rawReflectedSlice := make([]interface{}, 0)

	reflectedRawValues := reflect.ValueOf(rawInputValues)
	rawInputValuesKind := reflectedRawValues.Kind()
	switch rawInputValuesKind {
	case reflect.String:
		reflectedRawValuesStr := reflectedRawValues.String()
		rawSeparatedValues := strings.Split(reflectedRawValuesStr, ";")
		if len(rawSeparatedValues) <= 1 {
			rawSeparatedValues = strings.Split(reflectedRawValuesStr, ",")
		}

		for _, rawValue := range rawSeparatedValues {
			rawReflectedSlice = append(rawReflectedSlice, rawValue)
		}
	case reflect.Slice:
		for valueIndex := 0; valueIndex < reflectedRawValues.Len(); valueIndex++ {
			rawReflectedSlice = append(
				rawReflectedSlice, reflectedRawValues.Index(valueIndex).Interface(),
			)
		}
	default:
		rawReflectedSlice = append(rawReflectedSlice, rawInputValues)
	}

	for _, rawValue := range rawReflectedSlice {
		valueObject, err := valueObjectConstructor(rawValue)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("rawValue", rawValue))
			continue
		}

		parsedObjects = append(parsedObjects, valueObject)
	}

	return parsedObjects
}
