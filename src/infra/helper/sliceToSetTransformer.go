package infraHelper

func SliceToSetTransformer[TypedObject comparable](
	inputSlice []TypedObject,
) map[TypedObject]struct{} {
	set := map[TypedObject]struct{}{}
	for _, sliceKey := range inputSlice {
		set[sliceKey] = struct{}{}
	}
	return set
}
