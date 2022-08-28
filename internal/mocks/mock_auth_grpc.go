// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ports/auth_grpc.go

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	models "gitlab.com/g6834/team31/analytics/internal/domain/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockClientAuth is a mock of ClientAuth interface.
type MockClientAuth struct {
	ctrl     *gomock.Controller
	recorder *MockClientAuthMockRecorder
}

// MockClientAuthMockRecorder is the mock recorder for MockClientAuth.
type MockClientAuthMockRecorder struct {
	mock *MockClientAuth
}

// NewMockClientAuth creates a new mock instance.
func NewMockClientAuth(ctrl *gomock.Controller) *MockClientAuth {
	mock := &MockClientAuth{ctrl: ctrl}
	mock.recorder = &MockClientAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientAuth) EXPECT() *MockClientAuthMockRecorder {
	return m.recorder
}

// Validate mocks base method.
func (m *MockClientAuth) Validate(ctx context.Context, in models.JWTTokens) (models.ValidateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", ctx, in)
	ret0, _ := ret[0].(models.ValidateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockClientAuthMockRecorder) Validate(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockClientAuth)(nil).Validate), ctx, in)
}