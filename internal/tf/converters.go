package tf

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

func ExpandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func ExpandStringSet(configured *schema.Set) []string {
	return ExpandStringList(configured.List())
}
