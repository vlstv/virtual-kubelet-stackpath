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

// V1WatcherState v1 watcher state
//
// swagger:model v1WatcherState
type V1WatcherState string

func NewV1WatcherState(value V1WatcherState) *V1WatcherState {
	return &value
}

// Pointer returns a pointer to a freshly-allocated V1WatcherState.
func (m V1WatcherState) Pointer() *V1WatcherState {
	return &m
}

const (

	// V1WatcherStateWATCHERSTATEUNSPECIFIED captures enum value "WATCHER_STATE_UNSPECIFIED"
	V1WatcherStateWATCHERSTATEUNSPECIFIED V1WatcherState = "WATCHER_STATE_UNSPECIFIED"

	// V1WatcherStateLISTING captures enum value "LISTING"
	V1WatcherStateLISTING V1WatcherState = "LISTING"

	// V1WatcherStateWATCHING captures enum value "WATCHING"
	V1WatcherStateWATCHING V1WatcherState = "WATCHING"
)

// for schema
var v1WatcherStateEnum []interface{}

func init() {
	var res []V1WatcherState
	if err := json.Unmarshal([]byte(`["WATCHER_STATE_UNSPECIFIED","LISTING","WATCHING"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		v1WatcherStateEnum = append(v1WatcherStateEnum, v)
	}
}

func (m V1WatcherState) validateV1WatcherStateEnum(path, location string, value V1WatcherState) error {
	if err := validate.EnumCase(path, location, value, v1WatcherStateEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this v1 watcher state
func (m V1WatcherState) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateV1WatcherStateEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this v1 watcher state based on context it is used
func (m V1WatcherState) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}