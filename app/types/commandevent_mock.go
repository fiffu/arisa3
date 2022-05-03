// Code generated by MockGen. DO NOT EDIT.
// Source: commandevent.go

// Package types is a generated GoMock package.
package types

import (
	reflect "reflect"

	discordgo "github.com/bwmarrin/discordgo"
	gomock "github.com/golang/mock/gomock"
)

// MockICommandEvent is a mock of ICommandEvent interface.
type MockICommandEvent struct {
	ctrl     *gomock.Controller
	recorder *MockICommandEventMockRecorder
}

// MockICommandEventMockRecorder is the mock recorder for MockICommandEvent.
type MockICommandEventMockRecorder struct {
	mock *MockICommandEvent
}

// NewMockICommandEvent creates a new mock instance.
func NewMockICommandEvent(ctrl *gomock.Controller) *MockICommandEvent {
	mock := &MockICommandEvent{ctrl: ctrl}
	mock.recorder = &MockICommandEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICommandEvent) EXPECT() *MockICommandEventMockRecorder {
	return m.recorder
}

// Args mocks base method.
func (m *MockICommandEvent) Args() IArgs {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Args")
	ret0, _ := ret[0].(IArgs)
	return ret0
}

// Args indicates an expected call of Args.
func (mr *MockICommandEventMockRecorder) Args() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Args", reflect.TypeOf((*MockICommandEvent)(nil).Args))
}

// Command mocks base method.
func (m *MockICommandEvent) Command() ICommand {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Command")
	ret0, _ := ret[0].(ICommand)
	return ret0
}

// Command indicates an expected call of Command.
func (mr *MockICommandEventMockRecorder) Command() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Command", reflect.TypeOf((*MockICommandEvent)(nil).Command))
}

// Interaction mocks base method.
func (m *MockICommandEvent) Interaction() *discordgo.InteractionCreate {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Interaction")
	ret0, _ := ret[0].(*discordgo.InteractionCreate)
	return ret0
}

// Interaction indicates an expected call of Interaction.
func (mr *MockICommandEventMockRecorder) Interaction() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Interaction", reflect.TypeOf((*MockICommandEvent)(nil).Interaction))
}

// Respond mocks base method.
func (m *MockICommandEvent) Respond(arg0 ICommandResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Respond", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Respond indicates an expected call of Respond.
func (mr *MockICommandEventMockRecorder) Respond(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Respond", reflect.TypeOf((*MockICommandEvent)(nil).Respond), arg0)
}

// Session mocks base method.
func (m *MockICommandEvent) Session() *discordgo.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Session")
	ret0, _ := ret[0].(*discordgo.Session)
	return ret0
}

// Session indicates an expected call of Session.
func (mr *MockICommandEventMockRecorder) Session() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Session", reflect.TypeOf((*MockICommandEvent)(nil).Session))
}

// User mocks base method.
func (m *MockICommandEvent) User() *discordgo.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User")
	ret0, _ := ret[0].(*discordgo.User)
	return ret0
}

// User indicates an expected call of User.
func (mr *MockICommandEventMockRecorder) User() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockICommandEvent)(nil).User))
}
