package cursors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func TestNewCursorFromRequestLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		query string
		want  int
	}{
		{name: "missing falls back to the default", query: "", want: 20},
		{name: "valid value is kept", query: "limit=24", want: 24},
		{name: "zero falls back to the default", query: "limit=0", want: 20},
		{name: "negative falls back to the default", query: "limit=-5", want: 20},
		{name: "not a number falls back to the default", query: "limit=abc", want: 20},
		{name: "above the maximum is clamped", query: "limit=101", want: 100},
		{name: "int32 overflow is clamped", query: "limit=4294967297", want: 100},
		{name: "int64 overflow uses the default", query: "limit=99999999999999999999", want: 20},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequestWithContext(
				t.Context(),
				http.MethodGet,
				"/profiles?"+testCase.query,
				nil,
			)

			got := cursors.NewCursorFromRequest(request).Limit
			if got != testCase.want {
				t.Errorf("limit for %q = %d, want %d", testCase.query, got, testCase.want)
			}
		})
	}
}

func TestNewCursorClampsLimit(t *testing.T) {
	t.Parallel()

	if got := cursors.NewCursor(5000, nil).Limit; got != 100 {
		t.Errorf("NewCursor(5000).Limit = %d, want 100", got)
	}
}
