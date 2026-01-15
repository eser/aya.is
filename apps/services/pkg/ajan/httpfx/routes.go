package httpfx

import (
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx/uris"
)

//go:generate go tool stringer -type RouteParameterType -trimprefix RouteParameterType
type RouteParameterType int

const (
	RouteParameterTypeHeader RouteParameterType = iota
	RouteParameterTypeQuery
	RouteParameterTypePath
	RouteParameterTypeBody
)

type RouterParameterValidator func(inputString string) (string, error)

type RouterParameter struct {
	Name        string
	Description string
	Validators  []RouterParameterValidator
	Type        RouteParameterType
	IsRequired  bool
}

type RouteOpenAPISpecRequest struct {
	Model any
}

type RouteOpenAPISpecResponse struct {
	Model      any
	StatusCode int
	HasModel   bool
}

type RouteOpenAPISpec struct {
	OperationID string
	Summary     string
	Description string
	Tags        []string

	Requests   []RouteOpenAPISpecRequest
	Responses  []RouteOpenAPISpecResponse
	Deprecated bool
}

type Route struct {
	Pattern        *uris.Pattern
	Parameters     []RouterParameter
	Handlers       []Handler
	MuxHandlerFunc func(http.ResponseWriter, *http.Request)

	Spec RouteOpenAPISpec

	// frozen indicates this route can no longer be modified (set when router is frozen)
	frozen bool
}

func (r *Route) HasOperationID(operationID string) *Route {
	r.Spec.OperationID = operationID

	return r
}

func (r *Route) HasSummary(summary string) *Route {
	r.Spec.Summary = summary

	return r
}

func (r *Route) HasDescription(description string) *Route {
	r.Spec.Description = description

	return r
}

func (r *Route) HasTags(tags ...string) *Route {
	r.Spec.Tags = tags

	return r
}

func (r *Route) IsDeprecated() *Route {
	r.Spec.Deprecated = true

	return r
}

func (r *Route) HasPathParameter(name string, description string) *Route {
	r.Parameters = append(r.Parameters, RouterParameter{
		Type:        RouteParameterTypePath,
		Name:        name,
		Description: description,
		IsRequired:  true,

		Validators: []RouterParameterValidator{
			// func(inputString string) (string, error) {
			// 	return inputString, nil
			// },
		},
	})

	return r
}

func (r *Route) HasQueryParameter(name string, description string) *Route {
	r.Parameters = append(r.Parameters, RouterParameter{
		Type:        RouteParameterTypeQuery,
		Name:        name,
		Description: description,
		IsRequired:  true,

		Validators: []RouterParameterValidator{
			// func(inputString string) (string, error) {
			// 	return inputString, nil
			// },
		},
	})

	return r
}

func (r *Route) HasRequestModel(model any) *Route {
	r.Spec.Requests = append(r.Spec.Requests, RouteOpenAPISpecRequest{
		Model: model,
	})

	return r
}

func (r *Route) HasResponse(statusCode int) *Route {
	r.Spec.Responses = append(r.Spec.Responses, RouteOpenAPISpecResponse{
		StatusCode: statusCode,
		HasModel:   false,
		Model:      nil,
	})

	return r
}

func (r *Route) HasResponseModel(statusCode int, model any) *Route {
	r.Spec.Responses = append(r.Spec.Responses, RouteOpenAPISpecResponse{
		StatusCode: statusCode,
		HasModel:   true,
		Model:      model,
	})

	return r
}
