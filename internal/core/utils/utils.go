package utils

import "fmt"

func ArgInt(named map[string]any, key string, required bool) (int, error) {
	val, exists := named[key]
	if !exists {
		if required {
			return 0, fmt.Errorf("required parameter '%s' not found", key)
		}
		return 0, nil
	}

	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("parameter '%s' must be a number, got %T", key, val)
	}
}

func ArgIntSlice(named map[string]any, key string, required bool) ([]int, error) {
	val, exists := named[key]
	if !exists {
		if required {
			return nil, fmt.Errorf("required parameter '%s' not found", key)
		}
		return nil, nil
	}

	switch v := val.(type) {
	case []int:
		return v, nil
	case []any:
		result := make([]int, len(v))
		for i, item := range v {
			switch n := item.(type) {
			case int:
				result[i] = n
			case int64:
				result[i] = int(n)
			case float64:
				result[i] = int(n)
			default:
				return nil, fmt.Errorf("values[%d] must be a number, got %T", i, item)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("parameter '%s' must be a slice of numbers, got %T", key, val)
	}
}
