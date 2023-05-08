// Code generated by go-swagger; DO NOT EDIT.

package workload_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// V1GetMetricsRequestType The types of metrics that can be retrieved
//
// - BANDWIDTH: Bandwidth ingress and egress metrics
//   - INSTANCE: Instance CPU, memory, and filesystem metrics
//
// swagger:model v1GetMetricsRequestType
type V1GetMetricsRequestType string

func NewV1GetMetricsRequestType(value V1GetMetricsRequestType) *V1GetMetricsRequestType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated V1GetMetricsRequestType.
func (m V1GetMetricsRequestType) Pointer() *V1GetMetricsRequestType {
	return &m
}

const (

	// V1GetMetricsRequestTypeBANDWIDTH captures enum value "BANDWIDTH"
	V1GetMetricsRequestTypeBANDWIDTH V1GetMetricsRequestType = "BANDWIDTH"

	// V1GetMetricsRequestTypeINSTANCE captures enum value "INSTANCE"
	V1GetMetricsRequestTypeINSTANCE V1GetMetricsRequestType = "INSTANCE"
)

// for schema
var v1GetMetricsRequestTypeEnum []interface{}

func init() {
	var res []V1GetMetricsRequestType
	if err := json.Unmarshal([]byte(`["BANDWIDTH","INSTANCE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		v1GetMetricsRequestTypeEnum = append(v1GetMetricsRequestTypeEnum, v)
	}
}

func (m V1GetMetricsRequestType) validateV1GetMetricsRequestTypeEnum(path, location string, value V1GetMetricsRequestType) error {
	if err := validate.EnumCase(path, location, value, v1GetMetricsRequestTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this v1 get metrics request type
func (m V1GetMetricsRequestType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateV1GetMetricsRequestTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this v1 get metrics request type based on context it is used
func (m V1GetMetricsRequestType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}