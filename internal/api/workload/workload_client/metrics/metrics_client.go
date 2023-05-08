// Code generated by go-swagger; DO NOT EDIT.

package metrics

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new metrics API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for metrics API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	GetMetrics(params *GetMetricsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetMetricsOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
GetMetrics retrieves usage metrics for all workloads in a stack a specific workload or a specific instance in a workload
*/
func (a *Client) GetMetrics(params *GetMetricsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetMetricsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMetricsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetMetrics",
		Method:             "GET",
		PathPattern:        "/workload/v1/stacks/{stack_id}/metrics",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetMetricsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetMetricsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetMetricsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}