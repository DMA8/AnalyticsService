// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ports/task_grpc.go

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	models "gitlab.com/g6834/team31/analytics/internal/domain/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockServerTask is a mock of ServerTask interface.
type MockServerTask struct {
	ctrl     *gomock.Controller
	recorder *MockServerTaskMockRecorder
}

// MockServerTaskMockRecorder is the mock recorder for MockServerTask.
type MockServerTaskMockRecorder struct {
	mock *MockServerTask
}

// NewMockServerTask creates a new mock instance.
func NewMockServerTask(ctrl *gomock.Controller) *MockServerTask {
	mock := &MockServerTask{ctrl: ctrl}
	mock.recorder = &MockServerTaskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerTask) EXPECT() *MockServerTaskMockRecorder {
	return m.recorder
}

// PushMail mocks base method.
func (m *MockServerTask) PushMail(ctx context.Context, mail models.Mail) (models.TaskResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushMail", ctx, mail)
	ret0, _ := ret[0].(models.TaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushMail indicates an expected call of PushMail.
func (mr *MockServerTaskMockRecorder) PushMail(ctx, mail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushMail", reflect.TypeOf((*MockServerTask)(nil).PushMail), ctx, mail)
}

// PushTask mocks base method.
func (m *MockServerTask) PushTask(ctx context.Context, task models.Task, action, kind int) (models.TaskResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushTask", ctx, task, action, kind)
	ret0, _ := ret[0].(models.TaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushTask indicates an expected call of PushTask.
func (mr *MockServerTaskMockRecorder) PushTask(ctx, task, action, kind interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushTask", reflect.TypeOf((*MockServerTask)(nil).PushTask), ctx, task, action, kind)
}
