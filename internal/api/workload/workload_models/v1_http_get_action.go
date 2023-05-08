// Code generated by go-swagger; DO NOT EDIT.

package workload_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1HTTPGetAction Execute an HTTP GET request against an endpoint running on an instance
//
// swagger:model v1HTTPGetAction
type V1HTTPGetAction struct {

	// http headers
	HTTPHeaders V1StringMapEntry `json:"httpHeaders,omitempty"`

	// The path portion of the URL to request
	Path string `json:"path,omitempty"`

	// The TCP port to query in the HTTP request
	Port int32 `json:"port,omitempty"`

	// HTTP scheme to use in the HTTP request
	Scheme string `json:"scheme,omitempty"`
}

// Validate validates this v1 HTTP get action
func (m *V1HTTPGetAction) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHTTPHeaders(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1HTTPGetAction) validateHTTPHeaders(formats strfmt.Registry) error {
	if swag.IsZero(m.HTTPHeaders) { // not required
		return nil
	}

	if m.HTTPHeaders != nil {
		if err := m.HTTPHeaders.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("httpHeaders")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("httpHeaders")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this v1 HTTP get action based on the context it is used
func (m *V1HTTPGetAction) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateHTTPHeaders(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1HTTPGetAction) contextValidateHTTPHeaders(ctx context.Context, formats strfmt.Registry) error {

	if err := m.HTTPHeaders.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("httpHeaders")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("httpHeaders")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1HTTPGetAction) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1HTTPGetAction) UnmarshalBinary(b []byte) error {
	var res V1HTTPGetAction
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}