package lib

import (
	"strings"
)

func CaseInsensitiveSet(m *map[string]any, key string, value any) { //nolint:varnamelen
	for k := range *m {
		if strings.EqualFold(k, key) {
			(*m)[k] = value

			return
		}
	}

	(*m)[key] = value
}

// CaseInsensitiveGet retrieves a value from the map using case-insensitive key matching.
// Returns the value and true if found, or nil and false if not found.
func CaseInsensitiveGet(m *map[string]any, key string) (any, bool) {
	for k, v := range *m {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}

	return nil, false
}
