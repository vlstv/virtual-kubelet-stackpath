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

// V1InstanceConditionType v1 instance condition type
//
// swagger:model v1InstanceConditionType
type V1InstanceConditionType string

func NewV1InstanceConditionType(value V1InstanceConditionType) *V1InstanceConditionType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated V1InstanceConditionType.
func (m V1InstanceConditionType) Pointer() *V1InstanceConditionType {
	return &m
}

const (

	// V1InstanceConditionTypeREADY captures enum value "READY"
	V1InstanceConditionTypeREADY V1InstanceConditionType = "READY"
)

// for schema
var v1InstanceConditionTypeEnum []interface{}

func init() {
	var res []V1InstanceConditionType
	if err := json.Unmarshal([]byte(`["READY"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		v1InstanceConditionTypeEnum = append(v1InstanceConditionTypeEnum, v)
	}
}

func (m V1InstanceConditionType) validateV1InstanceConditionTypeEnum(path, location string, value V1InstanceConditionType) error {
	if err := validate.EnumCase(path, location, value, v1InstanceConditionTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this v1 instance condition type
func (m V1InstanceConditionType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateV1InstanceConditionTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this v1 instance condition type based on context it is used
func (m V1InstanceConditionType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}