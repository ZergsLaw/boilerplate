// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UpdatePassword update password
//
// swagger:model UpdatePassword
type UpdatePassword struct {

	// new
	// Required: true
	// Format: password
	New Password `json:"new"`

	// old
	// Required: true
	// Format: password
	Old Password `json:"old"`
}

// Validate validates this update password
func (m *UpdatePassword) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateNew(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOld(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UpdatePassword) validateNew(formats strfmt.Registry) error {

	if err := m.New.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("new")
		}
		return err
	}

	return nil
}

func (m *UpdatePassword) validateOld(formats strfmt.Registry) error {

	if err := m.Old.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("old")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *UpdatePassword) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UpdatePassword) UnmarshalBinary(b []byte) error {
	var res UpdatePassword
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
