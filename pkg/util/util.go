package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ToJsonString marshals the given data into a JSON string.
func ToJsonString(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// IncludesJson checks if all key-value pairs in the desired JSON string are present in the current JSON string.
func IncludesJson(currentJson, desiredJson string) (bool, error) {
	var current, desired map[string]interface{}

	// Parse the current JSON string
	if err := json.Unmarshal([]byte(currentJson), &current); err != nil {
		return false, fmt.Errorf("error parsing current JSON: %w", err)
	}

	// Parse the desired JSON string
	if err := json.Unmarshal([]byte(desiredJson), &desired); err != nil {
		return false, fmt.Errorf("error parsing desired JSON: %w", err)
	}

	// Check if all key-value pairs in desired are present in current
	return includes(current, desired), nil
}

// includes recursively checks if all key-value pairs in the desired map are present in the current map.
func includes(current, desired map[string]interface{}) bool {
	for key, desiredValue := range desired {
		currentValue, exists := current[key]
		if !exists {
			return false
		}

		if !reflect.DeepEqual(desiredValue, currentValue) {
			switch desiredValueTyped := desiredValue.(type) {
			case map[string]interface{}:
				currentValueTyped, ok := currentValue.(map[string]interface{})
				if !ok || !includes(currentValueTyped, desiredValueTyped) {
					return false
				}
			case []interface{}:
				currentValueTyped, ok := currentValue.([]interface{})
				if !ok || !sliceIncludes(currentValueTyped, desiredValueTyped) {
					return false
				}
			default:
				return false
			}
		}
	}
	return true
}

// sliceIncludes checks if all elements of the desired slice are present in the current slice.
func sliceIncludes(current, desired []interface{}) bool {
	for _, desiredItem := range desired {
		found := false
		for _, currentItem := range current {
			if reflect.DeepEqual(desiredItem, currentItem) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
