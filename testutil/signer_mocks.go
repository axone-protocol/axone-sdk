// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hyperledger/aries-framework-go/pkg/doc/verifiable (interfaces: Signer)
//
// Generated by this command:
//
//	mockgen -package testutil -destination testutil/signer_mocks.go github.com/hyperledger/aries-framework-go/pkg/doc/verifiable Signer
//

// Package testutil is a generated GoMock package.
package testutil

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSigner is a mock of Signer interface.
type MockSigner struct {
	ctrl     *gomock.Controller
	recorder *MockSignerMockRecorder
}

// MockSignerMockRecorder is the mock recorder for MockSigner.
type MockSignerMockRecorder struct {
	mock *MockSigner
}

// NewMockSigner creates a new mock instance.
func NewMockSigner(ctrl *gomock.Controller) *MockSigner {
	mock := &MockSigner{ctrl: ctrl}
	mock.recorder = &MockSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSigner) EXPECT() *MockSignerMockRecorder {
	return m.recorder
}

// Alg mocks base method.
func (m *MockSigner) Alg() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Alg")
	ret0, _ := ret[0].(string)
	return ret0
}

// Alg indicates an expected call of Alg.
func (mr *MockSignerMockRecorder) Alg() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Alg", reflect.TypeOf((*MockSigner)(nil).Alg))
}

// Sign mocks base method.
func (m *MockSigner) Sign(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockSignerMockRecorder) Sign(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockSigner)(nil).Sign), arg0)
}
