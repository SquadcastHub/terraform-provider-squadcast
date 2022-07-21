package tf

import (
	"github.com/mitchellh/mapstructure"
)

const EncoderStructTag = "tf"

type StateEncoder interface {
	Encode() (map[string]any, error)
}

func Encode(input any) (map[string]any, error) {
	var m map[string]any

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:               &m,
		TagName:              EncoderStructTag,
		IgnoreUntaggedFields: true,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(input)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func EncodeSlice[T StateEncoder](input []T) ([]any, error) {
	slice := make([]any, len(input))
	for i, v := range input {
		m, err := v.Encode()
		if err != nil {
			return nil, err
		}
		slice[i] = m
	}

	return slice, nil
}
