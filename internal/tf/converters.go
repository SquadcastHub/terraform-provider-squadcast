package tf

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func ExtractData(d *schema.ResourceData, key string) (map[string]interface{}, error) {
	data := d.Get(key)
	dataList, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("data for key %s is not a list", key)
	}

	if len(dataList) > 0 {
		dataMap, ok := dataList[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("data for key %s is not a map", key)
		}
		return dataMap, nil
	}

	return nil, nil
}
