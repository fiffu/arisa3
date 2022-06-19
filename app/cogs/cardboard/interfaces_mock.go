// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package cardboard is a generated GoMock package.
package cardboard

import (
	reflect "reflect"

	api "github.com/fiffu/arisa3/app/cogs/cardboard/api"
	types "github.com/fiffu/arisa3/app/types"
	gomock "github.com/golang/mock/gomock"
)

// MockIDomain is a mock of IDomain interface.
type MockIDomain struct {
	ctrl     *gomock.Controller
	recorder *MockIDomainMockRecorder
}

// MockIDomainMockRecorder is the mock recorder for MockIDomain.
type MockIDomainMockRecorder struct {
	mock *MockIDomain
}

// NewMockIDomain creates a new mock instance.
func NewMockIDomain(ctrl *gomock.Controller) *MockIDomain {
	mock := &MockIDomain{ctrl: ctrl}
	mock.recorder = &MockIDomainMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDomain) EXPECT() *MockIDomainMockRecorder {
	return m.recorder
}

// AliasTag mocks base method.
func (m *MockIDomain) AliasTag(guildID, actual, alias string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AliasTag", guildID, actual, alias)
	ret0, _ := ret[0].(error)
	return ret0
}

// AliasTag indicates an expected call of AliasTag.
func (mr *MockIDomainMockRecorder) AliasTag(guildID, actual, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AliasTag", reflect.TypeOf((*MockIDomain)(nil).AliasTag), guildID, actual, alias)
}

// DemoteTag mocks base method.
func (m *MockIDomain) DemoteTag(guildID, tagName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DemoteTag", guildID, tagName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DemoteTag indicates an expected call of DemoteTag.
func (mr *MockIDomainMockRecorder) DemoteTag(guildID, tagName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DemoteTag", reflect.TypeOf((*MockIDomain)(nil).DemoteTag), guildID, tagName)
}

// OmitTag mocks base method.
func (m *MockIDomain) OmitTag(guildID, tagName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OmitTag", guildID, tagName)
	ret0, _ := ret[0].(error)
	return ret0
}

// OmitTag indicates an expected call of OmitTag.
func (mr *MockIDomainMockRecorder) OmitTag(guildID, tagName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OmitTag", reflect.TypeOf((*MockIDomain)(nil).OmitTag), guildID, tagName)
}

// PostsResult mocks base method.
func (m *MockIDomain) PostsResult(arg0 IQueryPosts, arg1 []*api.Post) (types.IEmbed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostsResult", arg0, arg1)
	ret0, _ := ret[0].(types.IEmbed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostsResult indicates an expected call of PostsResult.
func (mr *MockIDomainMockRecorder) PostsResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostsResult", reflect.TypeOf((*MockIDomain)(nil).PostsResult), arg0, arg1)
}

// PostsSearch mocks base method.
func (m *MockIDomain) PostsSearch(arg0 IQueryPosts) ([]*api.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostsSearch", arg0)
	ret0, _ := ret[0].([]*api.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostsSearch indicates an expected call of PostsSearch.
func (mr *MockIDomainMockRecorder) PostsSearch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostsSearch", reflect.TypeOf((*MockIDomain)(nil).PostsSearch), arg0)
}

// PromoteTag mocks base method.
func (m *MockIDomain) PromoteTag(guildID, tagName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PromoteTag", guildID, tagName)
	ret0, _ := ret[0].(error)
	return ret0
}

// PromoteTag indicates an expected call of PromoteTag.
func (mr *MockIDomainMockRecorder) PromoteTag(guildID, tagName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PromoteTag", reflect.TypeOf((*MockIDomain)(nil).PromoteTag), guildID, tagName)
}

// MockIQueryPosts is a mock of IQueryPosts interface.
type MockIQueryPosts struct {
	ctrl     *gomock.Controller
	recorder *MockIQueryPostsMockRecorder
}

// MockIQueryPostsMockRecorder is the mock recorder for MockIQueryPosts.
type MockIQueryPostsMockRecorder struct {
	mock *MockIQueryPosts
}

// NewMockIQueryPosts creates a new mock instance.
func NewMockIQueryPosts(ctrl *gomock.Controller) *MockIQueryPosts {
	mock := &MockIQueryPosts{ctrl: ctrl}
	mock.recorder = &MockIQueryPostsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIQueryPosts) EXPECT() *MockIQueryPostsMockRecorder {
	return m.recorder
}

// GuildID mocks base method.
func (m *MockIQueryPosts) GuildID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuildID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GuildID indicates an expected call of GuildID.
func (mr *MockIQueryPostsMockRecorder) GuildID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuildID", reflect.TypeOf((*MockIQueryPosts)(nil).GuildID))
}

// MagicMode mocks base method.
func (m *MockIQueryPosts) MagicMode() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MagicMode")
	ret0, _ := ret[0].(bool)
	return ret0
}

// MagicMode indicates an expected call of MagicMode.
func (mr *MockIQueryPostsMockRecorder) MagicMode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MagicMode", reflect.TypeOf((*MockIQueryPosts)(nil).MagicMode))
}

// SetTerm mocks base method.
func (m *MockIQueryPosts) SetTerm(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTerm", arg0)
}

// SetTerm indicates an expected call of SetTerm.
func (mr *MockIQueryPostsMockRecorder) SetTerm(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTerm", reflect.TypeOf((*MockIQueryPosts)(nil).SetTerm), arg0)
}

// Tags mocks base method.
func (m *MockIQueryPosts) Tags() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tags")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Tags indicates an expected call of Tags.
func (mr *MockIQueryPostsMockRecorder) Tags() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tags", reflect.TypeOf((*MockIQueryPosts)(nil).Tags))
}

// Term mocks base method.
func (m *MockIQueryPosts) Term() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Term")
	ret0, _ := ret[0].(string)
	return ret0
}

// Term indicates an expected call of Term.
func (mr *MockIQueryPostsMockRecorder) Term() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Term", reflect.TypeOf((*MockIQueryPosts)(nil).Term))
}

// MockIRepository is a mock of IRepository interface.
type MockIRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRepositoryMockRecorder
}

// MockIRepositoryMockRecorder is the mock recorder for MockIRepository.
type MockIRepositoryMockRecorder struct {
	mock *MockIRepository
}

// NewMockIRepository creates a new mock instance.
func NewMockIRepository(ctrl *gomock.Controller) *MockIRepository {
	mock := &MockIRepository{ctrl: ctrl}
	mock.recorder = &MockIRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRepository) EXPECT() *MockIRepositoryMockRecorder {
	return m.recorder
}

// GetAliases mocks base method.
func (m *MockIRepository) GetAliases(guildID string) (map[Alias]Actual, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAliases", guildID)
	ret0, _ := ret[0].(map[Alias]Actual)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAliases indicates an expected call of GetAliases.
func (mr *MockIRepositoryMockRecorder) GetAliases(guildID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAliases", reflect.TypeOf((*MockIRepository)(nil).GetAliases), guildID)
}

// GetDemotes mocks base method.
func (m *MockIRepository) GetDemotes(guildID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDemotes", guildID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDemotes indicates an expected call of GetDemotes.
func (mr *MockIRepositoryMockRecorder) GetDemotes(guildID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDemotes", reflect.TypeOf((*MockIRepository)(nil).GetDemotes), guildID)
}

// GetOmits mocks base method.
func (m *MockIRepository) GetOmits(guildID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOmits", guildID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOmits indicates an expected call of GetOmits.
func (mr *MockIRepositoryMockRecorder) GetOmits(guildID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOmits", reflect.TypeOf((*MockIRepository)(nil).GetOmits), guildID)
}

// GetPromotes mocks base method.
func (m *MockIRepository) GetPromotes(guildID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPromotes", guildID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPromotes indicates an expected call of GetPromotes.
func (mr *MockIRepositoryMockRecorder) GetPromotes(guildID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPromotes", reflect.TypeOf((*MockIRepository)(nil).GetPromotes), guildID)
}

// GetTagOperations mocks base method.
func (m *MockIRepository) GetTagOperations(guildID string) (map[string]TagOperation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagOperations", guildID)
	ret0, _ := ret[0].(map[string]TagOperation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagOperations indicates an expected call of GetTagOperations.
func (mr *MockIRepositoryMockRecorder) GetTagOperations(guildID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagOperations", reflect.TypeOf((*MockIRepository)(nil).GetTagOperations), guildID)
}

// SetAlias mocks base method.
func (m *MockIRepository) SetAlias(guildID string, ali Alias, act Actual) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAlias", guildID, ali, act)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAlias indicates an expected call of SetAlias.
func (mr *MockIRepositoryMockRecorder) SetAlias(guildID, ali, act interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAlias", reflect.TypeOf((*MockIRepository)(nil).SetAlias), guildID, ali, act)
}

// SetDemote mocks base method.
func (m *MockIRepository) SetDemote(guildID, tag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDemote", guildID, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDemote indicates an expected call of SetDemote.
func (mr *MockIRepositoryMockRecorder) SetDemote(guildID, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDemote", reflect.TypeOf((*MockIRepository)(nil).SetDemote), guildID, tag)
}

// SetOmit mocks base method.
func (m *MockIRepository) SetOmit(guildID, tag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOmit", guildID, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOmit indicates an expected call of SetOmit.
func (mr *MockIRepositoryMockRecorder) SetOmit(guildID, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOmit", reflect.TypeOf((*MockIRepository)(nil).SetOmit), guildID, tag)
}

// SetPromote mocks base method.
func (m *MockIRepository) SetPromote(guildID, tag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPromote", guildID, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPromote indicates an expected call of SetPromote.
func (mr *MockIRepositoryMockRecorder) SetPromote(guildID, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPromote", reflect.TypeOf((*MockIRepository)(nil).SetPromote), guildID, tag)
}
