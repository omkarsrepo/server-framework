package json

import (
	"fmt"
	"strings"
)

func GetObjectValue(data *map[string]interface{}, key string) (interface{}, error) {
	parts := strings.Split(key, ".")
	current := data

	for _, part := range parts[:len(parts)-1] {
		if currentMap, ok := (*current)[part].(map[string]interface{}); ok {
			current = &currentMap
		} else {
			return nil, fmt.Errorf("key %s not found or not a map", part)
		}
	}

	finalKey := parts[len(parts)-1]
	if value, ok := (*current)[finalKey]; ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("key %s not found", finalKey)
	}
}
