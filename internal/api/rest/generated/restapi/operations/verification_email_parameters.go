// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	models "github.com/zergslaw/users/internal/api/rest/generated/models"
)

// NewVerificationEmailParams creates a new VerificationEmailParams object
// no default values defined in spec.
func NewVerificationEmailParams() VerificationEmailParams {

	return VerificationEmailParams{}
}

// VerificationEmailParams contains all the bound params for the verification email operation
// typically these are obtained from a http.Request
//
// swagger:parameters verificationEmail
type VerificationEmailParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	Email models.Email
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewVerificationEmailParams() beforehand.
func (o *VerificationEmailParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.Email
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("email", "body"))
			} else {
				res = append(res, errors.NewParseError("email", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Email = body
			}
		}
	} else {
		res = append(res, errors.Required("email", "body"))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
