package util

import (
	"fmt"
)

func StringMapGetDefault(m map[string]string, k string, def string) string {
	v, ok := m[k]
	if !ok {
		v = def
	}
	return v
}

func StringFromJson(m map[string]interface{}, k string) (string, error) {

	sInterface, ok := m[k]
	if !ok {
		return "", fmt.Errorf("Failed to find key %v", k)
	}

	s, ok := sInterface.(string)
	if !ok {
		return "", fmt.Errorf("Failed to convert key %v to string", k)
	}

	return s, nil
}

func MapInterfaceFromJson(m map[string]interface{}, k string) (map[string]interface{}, error) {

	m2Interface, ok := m[k]
	if !ok {
		return nil, fmt.Errorf("Failed to find key %v", k)
	}

	m2, ok := m2Interface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Failed to convert key %v to string", k)
	}

	return m2, nil
}
