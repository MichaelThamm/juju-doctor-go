package utils

import (
	"fmt"

	"github.com/canonical/starlark/starlark"
)

// ToStarlarkDict converts a Go map to a Starlark dictionary.
func ToStarlarkDict(rawData map[string]interface{}) (starlark.Value, error) {
	starlarkDict := starlark.NewDict(len(rawData))
	for key, value := range rawData {
		starlarkValue, err := ToStarlarkValue(value)
		if err != nil {
			return nil, fmt.Errorf("error converting key %q: %v", key, err)
		}
		starlarkDict.SetKey(starlark.String(key), starlarkValue)
	}
	return starlarkDict, nil
}

// ToStarlarkValue converts Go types to Starlark values.
func ToStarlarkValue(value interface{}) (starlark.Value, error) {
	switch v := value.(type) {
	case string:
		return starlark.String(v), nil
	case int:
		return starlark.MakeInt(v), nil
	case int64:
		return starlark.MakeInt64(v), nil
	case float64:
		return starlark.Float(v), nil
	case bool:
		return starlark.Bool(v), nil
	case map[string]interface{}:
		return ToStarlarkDict(v)
	case []interface{}:
		starlarkList := make([]starlark.Value, len(v))
		for i, elem := range v {
			starlarkElem, err := ToStarlarkValue(elem)
			if err != nil {
				return nil, fmt.Errorf("error converting list element %d: %v", i, err)
			}
			starlarkList[i] = starlarkElem
		}
		return starlark.NewList(starlarkList), nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}
