package cursors

import (
	"net/http"
	"strconv"
	"strings"
)

func NewCursorFromRequest(r *http.Request) *Cursor {
	queryValues := r.URL.Query()

	limitStr := queryValues.Get("limit")
	limit := getLimitFromString(limitStr)

	offset := queryValues.Get("offset")

	sortStr := queryValues.Get("sort")
	sortBy, sortDir := getSortFromString(sortStr)

	filters := make(map[string]string)

	for key, values := range queryValues {
		if strings.HasPrefix(key, "filter_") && len(values) > 0 {
			filterKey := strings.TrimPrefix(key, "filter_")
			filters[filterKey] = values[0]
		}
	}

	return &Cursor{
		Limit:  limit,
		Offset: &offset,

		SortBy:  sortBy,
		SortDir: sortDir,

		Filters: filters,
	}
}

func getLimitFromString(str string) int {
	if str == "" {
		return defaultLimit
	}

	val, err := strconv.Atoi(str)
	if err != nil || val <= 0 {
		return defaultLimit
	}

	return val
}

func getSortFromString(str string) (string, string) {
	if str == "" {
		return defaultSortBy, defaultSortDir
	}

	sortParts := strings.SplitN(str, " ", 2) //nolint:mnd

	if len(sortParts) == 0 {
		return defaultSortBy, defaultSortDir
	}

	if len(sortParts) == 1 {
		return sortParts[0], defaultSortDir
	}

	return sortParts[0], sortParts[1]
}
