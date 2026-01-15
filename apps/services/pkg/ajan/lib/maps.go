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
