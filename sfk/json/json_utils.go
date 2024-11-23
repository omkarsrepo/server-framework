package json

import (
	"fmt"
	"strings"
)

func GetObjectValue(object *map[string]interface{}, keyInObject string) (interface{}, error) {
	parts := strings.Split(keyInObject, ".")
	current := object

	for _, part := range parts[:len(parts)-1] {
		if currentMap, ok := (*current)[part].(map[string]interface{}); ok {
			current = &currentMap
		} else {
			return nil, fmt.Errorf("object key '%s' not found or object is invalid", part)
		}
	}

	finalKey := parts[len(parts)-1]
	if value, ok := (*current)[finalKey]; ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("object key '%s' not found", finalKey)
	}
}
