// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
)

// NewV2GetClusterDefaultConfigParams creates a new V2GetClusterDefaultConfigParams object
// no default values defined in spec.
func NewV2GetClusterDefaultConfigParams() V2GetClusterDefaultConfigParams {

	return V2GetClusterDefaultConfigParams{}
}

// V2GetClusterDefaultConfigParams contains all the bound params for the v2 get cluster default config operation
// typically these are obtained from a http.Request
//
// swagger:parameters V2GetClusterDefaultConfig
type V2GetClusterDefaultConfigParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewV2GetClusterDefaultConfigParams() beforehand.
func (o *V2GetClusterDefaultConfigParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
