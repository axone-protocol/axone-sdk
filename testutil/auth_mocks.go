// Code generated by MockGen. DO NOT EDIT.
// Source: auth/proxy.go
//
// Generated by this command:
//
//	mockgen -source=auth/proxy.go -package testutil -destination testutil/auth_mocks.go
//

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	auth "github.com/axone-protocol/axone-sdk/auth"
	gomock "go.uber.org/mock/gomock"
)

// MockProxy is a mock of Proxy interface.
type MockProxy struct {
	ctrl     *gomock.Controller
	recorder *MockProxyMockRecorder
}

// MockProxyMockRecorder is the mock recorder for MockProxy.
type MockProxyMockRecorder struct {
	mock *MockProxy
}

// NewMockProxy creates a new mock instance.
func NewMockProxy(ctrl *gomock.Controller) *MockProxy {
	mock := &MockProxy{ctrl: ctrl}
	mock.recorder = &MockProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProxy) EXPECT() *MockProxyMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockProxy) Authenticate(ctx context.Context, credential []byte) (*auth.Identity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", ctx, credential)
	ret0, _ := ret[0].(*auth.Identity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockProxyMockRecorder) Authenticate(ctx, credential any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockProxy)(nil).Authenticate), ctx, credential)
}
