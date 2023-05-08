// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client/instance (interfaces: ClientService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	runtime "github.com/go-openapi/runtime"
	gomock "github.com/golang/mock/gomock"
	instance "github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client/instance"
)

// InstanceClientService is a mock of ClientService interface.
type InstanceClientService struct {
	ctrl     *gomock.Controller
	recorder *InstanceClientServiceMockRecorder
}

// InstanceClientServiceMockRecorder is the mock recorder for InstanceClientService.
type InstanceClientServiceMockRecorder struct {
	mock *InstanceClientService
}

// NewInstanceClientService creates a new mock instance.
func NewInstanceClientService(ctrl *gomock.Controller) *InstanceClientService {
	mock := &InstanceClientService{ctrl: ctrl}
	mock.recorder = &InstanceClientServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *InstanceClientService) EXPECT() *InstanceClientServiceMockRecorder {
	return m.recorder
}

// GetWorkloadInstance mocks base method.
func (m *InstanceClientService) GetWorkloadInstance(arg0 *instance.GetWorkloadInstanceParams, arg1 runtime.ClientAuthInfoWriter, arg2 ...instance.ClientOption) (*instance.GetWorkloadInstanceOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetWorkloadInstance", varargs...)
	ret0, _ := ret[0].(*instance.GetWorkloadInstanceOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkloadInstance indicates an expected call of GetWorkloadInstance.
func (mr *InstanceClientServiceMockRecorder) GetWorkloadInstance(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkloadInstance", reflect.TypeOf((*InstanceClientService)(nil).GetWorkloadInstance), varargs...)
}

// RestartInstance mocks base method.
func (m *InstanceClientService) RestartInstance(arg0 *instance.RestartInstanceParams, arg1 runtime.ClientAuthInfoWriter, arg2 ...instance.ClientOption) (*instance.RestartInstanceNoContent, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RestartInstance", varargs...)
	ret0, _ := ret[0].(*instance.RestartInstanceNoContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestartInstance indicates an expected call of RestartInstance.
func (mr *InstanceClientServiceMockRecorder) RestartInstance(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestartInstance", reflect.TypeOf((*InstanceClientService)(nil).RestartInstance), varargs...)
}

// SetTransport mocks base method.
func (m *InstanceClientService) SetTransport(arg0 runtime.ClientTransport) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTransport", arg0)
}

// SetTransport indicates an expected call of SetTransport.
func (mr *InstanceClientServiceMockRecorder) SetTransport(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTransport", reflect.TypeOf((*InstanceClientService)(nil).SetTransport), arg0)
}