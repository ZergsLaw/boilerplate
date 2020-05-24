// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/zergslaw/boilerplate/internal/api/web/generated/models"
)

// NewCreateUserParams creates a new CreateUserParams object
// with the default values initialized.
func NewCreateUserParams() *CreateUserParams {
	var ()
	return &CreateUserParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateUserParamsWithTimeout creates a new CreateUserParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateUserParamsWithTimeout(timeout time.Duration) *CreateUserParams {
	var ()
	return &CreateUserParams{

		timeout: timeout,
	}
}

// NewCreateUserParamsWithContext creates a new CreateUserParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateUserParamsWithContext(ctx context.Context) *CreateUserParams {
	var ()
	return &CreateUserParams{

		Context: ctx,
	}
}

// NewCreateUserParamsWithHTTPClient creates a new CreateUserParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateUserParamsWithHTTPClient(client *http.Client) *CreateUserParams {
	var ()
	return &CreateUserParams{
		HTTPClient: client,
	}
}

/*CreateUserParams contains all the parameters to send to the API endpoint
for the create user operation typically these are written to a http.Request
*/
type CreateUserParams struct {

	/*Args*/
	Args *models.CreateUserParams

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create user params
func (o *CreateUserParams) WithTimeout(timeout time.Duration) *CreateUserParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create user params
func (o *CreateUserParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create user params
func (o *CreateUserParams) WithContext(ctx context.Context) *CreateUserParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create user params
func (o *CreateUserParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create user params
func (o *CreateUserParams) WithHTTPClient(client *http.Client) *CreateUserParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create user params
func (o *CreateUserParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithArgs adds the args to the create user params
func (o *CreateUserParams) WithArgs(args *models.CreateUserParams) *CreateUserParams {
	o.SetArgs(args)
	return o
}

// SetArgs adds the args to the create user params
func (o *CreateUserParams) SetArgs(args *models.CreateUserParams) {
	o.Args = args
}

// WriteToRequest writes these params to a swagger request
func (o *CreateUserParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Args != nil {
		if err := r.SetBodyParam(o.Args); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
