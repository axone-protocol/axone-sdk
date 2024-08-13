// Code generated by MockGen. DO NOT EDIT.
// Source: credential/parser.go
//
// Generated by this command:
//
//	mockgen -source=credential/parser.go -package testutil -destination testutil/credential_mocks.go
//

// Package testutil is a generated GoMock package.
package testutil

import (
	reflect "reflect"

	credential "github.com/axone-protocol/axone-sdk/credential"
	verifiable "github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	gomock "go.uber.org/mock/gomock"
)

// MockClaim is a mock of Claim interface.
type MockClaim struct {
	ctrl     *gomock.Controller
	recorder *MockClaimMockRecorder
}

// MockClaimMockRecorder is the mock recorder for MockClaim.
type MockClaimMockRecorder struct {
	mock *MockClaim
}

// NewMockClaim creates a new mock instance.
func NewMockClaim(ctrl *gomock.Controller) *MockClaim {
	mock := &MockClaim{ctrl: ctrl}
	mock.recorder = &MockClaimMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClaim) EXPECT() *MockClaimMockRecorder {
	return m.recorder
}

// From mocks base method.
func (m *MockClaim) From(vc *verifiable.Credential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "From", vc)
	ret0, _ := ret[0].(error)
	return ret0
}

// From indicates an expected call of From.
func (mr *MockClaimMockRecorder) From(vc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "From", reflect.TypeOf((*MockClaim)(nil).From), vc)
}

// MockParser is a mock of Parser interface.
type MockParser[T credential.Claim] struct {
	ctrl     *gomock.Controller
	recorder *MockParserMockRecorder[T]
}

// MockParserMockRecorder is the mock recorder for MockParser.
type MockParserMockRecorder[T credential.Claim] struct {
	mock *MockParser[T]
}

// NewMockParser creates a new mock instance.
func NewMockParser[T credential.Claim](ctrl *gomock.Controller) *MockParser[T] {
	mock := &MockParser[T]{ctrl: ctrl}
	mock.recorder = &MockParserMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParser[T]) EXPECT() *MockParserMockRecorder[T] {
	return m.recorder
}

// ParseSigned mocks base method.
func (m *MockParser[T]) ParseSigned(raw []byte) (T, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseSigned", raw)
	ret0, _ := ret[0].(T)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseSigned indicates an expected call of ParseSigned.
func (mr *MockParserMockRecorder[T]) ParseSigned(raw any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseSigned", reflect.TypeOf((*MockParser[T])(nil).ParseSigned), raw)
}