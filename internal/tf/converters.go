package tf

type M = map[string]any

func List[T any](item T) []T {
	return []T{item}
}

func AssumeMaps(list any) []M {
	return list.([]M)
}

func ListToSlice[T any](list any) []T {
	islice := list.([]any)
	slicelen := len(islice)
	slice := make([]T, slicelen, slicelen)

	if slicelen == 0 {
		return slice
	}

	for i, v := range islice {
		slice[i] = v.(T)
	}

	return slice
}
