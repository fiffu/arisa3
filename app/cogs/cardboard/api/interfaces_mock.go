// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package api is a generated GoMock package.
package api

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIClient is a mock of IClient interface.
type MockIClient struct {
	ctrl     *gomock.Controller
	recorder *MockIClientMockRecorder
}

// MockIClientMockRecorder is the mock recorder for MockIClient.
type MockIClientMockRecorder struct {
	mock *MockIClient
}

// NewMockIClient creates a new mock instance.
func NewMockIClient(ctrl *gomock.Controller) *MockIClient {
	mock := &MockIClient{ctrl: ctrl}
	mock.recorder = &MockIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIClient) EXPECT() *MockIClientMockRecorder {
	return m.recorder
}

// GetPosts mocks base method.
func (m *MockIClient) GetPosts(tags []string) ([]*Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPosts", tags)
	ret0, _ := ret[0].([]*Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPosts indicates an expected call of GetPosts.
func (mr *MockIClientMockRecorder) GetPosts(tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPosts", reflect.TypeOf((*MockIClient)(nil).GetPosts), tags)
}

// GetTags mocks base method.
func (m *MockIClient) GetTags(tags []string) (map[string]*Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTags", tags)
	ret0, _ := ret[0].(map[string]*Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTags indicates an expected call of GetTags.
func (mr *MockIClientMockRecorder) GetTags(tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTags", reflect.TypeOf((*MockIClient)(nil).GetTags), tags)
}

// GetTagsMatching mocks base method.
func (m *MockIClient) GetTagsMatching(pattern string) ([]*Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagsMatching", pattern)
	ret0, _ := ret[0].([]*Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagsMatching indicates an expected call of GetTagsMatching.
func (mr *MockIClientMockRecorder) GetTagsMatching(pattern interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagsMatching", reflect.TypeOf((*MockIClient)(nil).GetTagsMatching), pattern)
}
