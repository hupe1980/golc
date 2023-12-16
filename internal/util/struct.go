package util

import (
	"github.com/mitchellh/mapstructure"
)

func StructToMap(obj any) map[string]any {
	result := map[string]any{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "map",
		Result:  &result,
	})
	if err != nil {
		panic(err)
	}

	if err := decoder.Decode(obj); err != nil {
		panic(err)
	}

	return result
}
