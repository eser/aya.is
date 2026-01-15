package httpfx

import (
	"github.com/eser/aya.is/services/pkg/ajan/results"
)

type Result struct { //nolint:errname
	InnerRedirectToURI string
	results.Result

	InnerBody []byte

	InnerStatusCode int
}

func (r Result) StatusCode() int {
	return r.InnerStatusCode
}

func (r Result) Body() []byte {
	return r.InnerBody
}

func (r Result) RedirectToURI() string {
	return r.InnerRedirectToURI
}
