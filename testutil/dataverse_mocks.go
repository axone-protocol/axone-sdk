// Code generated by MockGen. DO NOT EDIT.
// Source: dataverse/client.go
//
// Generated by this command:
//
//	mockgen -source=dataverse/client.go -package testutil -destination testutil/dataverse_mocks.go
//

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// ExecGov mocks base method.
func (m *MockClient) ExecGov(arg0 context.Context, arg1, arg2 string) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecGov", arg0, arg1, arg2)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecGov indicates an expected call of ExecGov.
func (mr *MockClientMockRecorder) ExecGov(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecGov", reflect.TypeOf((*MockClient)(nil).ExecGov), arg0, arg1, arg2)
}

// GetGovAddr mocks base method.
func (m *MockClient) GetGovAddr(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGovAddr", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGovAddr indicates an expected call of GetGovAddr.
func (mr *MockClientMockRecorder) GetGovAddr(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGovAddr", reflect.TypeOf((*MockClient)(nil).GetGovAddr), arg0)
}
