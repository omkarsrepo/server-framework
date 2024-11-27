package json

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
)

func extractValue[T any](object *map[string]interface{}, keyInObject string, castFunc func(interface{}) (T, bool)) (T, error) {
	var zero T
	parts := strings.Split(keyInObject, ".")
	current := object

	for _, part := range parts[:len(parts)-1] {
		if currentMap, ok := (*current)[part].(map[string]interface{}); ok {
			current = &currentMap
		} else {
			return zero, fmt.Errorf("object key '%s' not found or object is invalid", part)
		}
	}

	finalKey := parts[len(parts)-1]
	if value, ok := (*current)[finalKey]; ok {
		if result, ok := castFunc(value); ok {
			return result, nil
		}
		return zero, fmt.Errorf("value at key '%s' cannot be cast to the expected type", finalKey)
	} else {
		return zero, fmt.Errorf("object key '%s' not found", finalKey)
	}
}

func ValueOf[T constraints.Ordered](object *map[string]interface{}, keyInObject string) (T, error) {
	return extractValue[T](object, keyInObject, func(value interface{}) (T, bool) {
		result, ok := value.(T)
		return result, ok
	})
}

func ListValueOf[T constraints.Ordered](object *map[string]interface{}, keyInObject string) ([]T, error) {
	return extractValue[[]T](object, keyInObject, func(value interface{}) ([]T, bool) {
		result, ok := value.([]T)
		return result, ok
	})
}

func AnyValueOf[T any](object *map[string]interface{}, keyInObject string) (T, error) {
	return extractValue[T](object, keyInObject, func(value interface{}) (T, bool) {
		result, ok := value.(T)
		return result, ok
	})
}
