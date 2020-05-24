// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/zergslaw/boilerplate/internal/api/web/generated/models"
)

// UpdatePasswordReader is a Reader for the UpdatePassword structure.
type UpdatePasswordReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdatePasswordReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewUpdatePasswordNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewUpdatePasswordDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdatePasswordNoContent creates a UpdatePasswordNoContent with default headers values
func NewUpdatePasswordNoContent() *UpdatePasswordNoContent {
	return &UpdatePasswordNoContent{}
}

/*UpdatePasswordNoContent handles this case with default header values.

The server successfully processed the request and is not returning any content.
*/
type UpdatePasswordNoContent struct {
}

func (o *UpdatePasswordNoContent) Error() string {
	return fmt.Sprintf("[PATCH /user/password][%d] updatePasswordNoContent ", 204)
}

func (o *UpdatePasswordNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdatePasswordDefault creates a UpdatePasswordDefault with default headers values
func NewUpdatePasswordDefault(code int) *UpdatePasswordDefault {
	return &UpdatePasswordDefault{
		_statusCode: code,
	}
}

/*UpdatePasswordDefault handles this case with default header values.

Generic error response.
*/
type UpdatePasswordDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the update password default response
func (o *UpdatePasswordDefault) Code() int {
	return o._statusCode
}

func (o *UpdatePasswordDefault) Error() string {
	return fmt.Sprintf("[PATCH /user/password][%d] updatePassword default  %+v", o._statusCode, o.Payload)
}

func (o *UpdatePasswordDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *UpdatePasswordDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
