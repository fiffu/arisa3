// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package colours is a generated GoMock package.
package colours

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockIColoursDomain is a mock of IColoursDomain interface.
type MockIColoursDomain struct {
	ctrl     *gomock.Controller
	recorder *MockIColoursDomainMockRecorder
}

// MockIColoursDomainMockRecorder is the mock recorder for MockIColoursDomain.
type MockIColoursDomainMockRecorder struct {
	mock *MockIColoursDomain
}

// NewMockIColoursDomain creates a new mock instance.
func NewMockIColoursDomain(ctrl *gomock.Controller) *MockIColoursDomain {
	mock := &MockIColoursDomain{ctrl: ctrl}
	mock.recorder = &MockIColoursDomainMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIColoursDomain) EXPECT() *MockIColoursDomainMockRecorder {
	return m.recorder
}

// AssignColourRole mocks base method.
func (m *MockIColoursDomain) AssignColourRole(arg0 context.Context, arg1 IDomainSession, arg2 IDomainMember, arg3 IDomainRole) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignColourRole", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignColourRole indicates an expected call of AssignColourRole.
func (mr *MockIColoursDomainMockRecorder) AssignColourRole(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignColourRole", reflect.TypeOf((*MockIColoursDomain)(nil).AssignColourRole), arg0, arg1, arg2, arg3)
}

// CreateColourRole mocks base method.
func (m *MockIColoursDomain) CreateColourRole(arg0 context.Context, arg1 IDomainSession, arg2 IDomainMember, arg3 *Colour) (IDomainRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateColourRole", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(IDomainRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateColourRole indicates an expected call of CreateColourRole.
func (mr *MockIColoursDomainMockRecorder) CreateColourRole(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateColourRole", reflect.TypeOf((*MockIColoursDomain)(nil).CreateColourRole), arg0, arg1, arg2, arg3)
}

// Freeze mocks base method.
func (m *MockIColoursDomain) Freeze(arg0 context.Context, arg1 IDomainMember) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Freeze", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Freeze indicates an expected call of Freeze.
func (mr *MockIColoursDomainMockRecorder) Freeze(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Freeze", reflect.TypeOf((*MockIColoursDomain)(nil).Freeze), arg0, arg1)
}

// GetColourRole mocks base method.
func (m *MockIColoursDomain) GetColourRole(arg0 context.Context, arg1 IDomainMember) IDomainRole {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetColourRole", arg0, arg1)
	ret0, _ := ret[0].(IDomainRole)
	return ret0
}

// GetColourRole indicates an expected call of GetColourRole.
func (mr *MockIColoursDomainMockRecorder) GetColourRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetColourRole", reflect.TypeOf((*MockIColoursDomain)(nil).GetColourRole), arg0, arg1)
}

// GetColourRoleHeight mocks base method.
func (m *MockIColoursDomain) GetColourRoleHeight(arg0 context.Context, arg1 IDomainSession, arg2 IDomainGuild) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetColourRoleHeight", arg0, arg1, arg2)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetColourRoleHeight indicates an expected call of GetColourRoleHeight.
func (mr *MockIColoursDomainMockRecorder) GetColourRoleHeight(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetColourRoleHeight", reflect.TypeOf((*MockIColoursDomain)(nil).GetColourRoleHeight), arg0, arg1, arg2)
}

// GetColourRoleName mocks base method.
func (m *MockIColoursDomain) GetColourRoleName(arg0 context.Context, arg1 IDomainMember) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetColourRoleName", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetColourRoleName indicates an expected call of GetColourRoleName.
func (mr *MockIColoursDomainMockRecorder) GetColourRoleName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetColourRoleName", reflect.TypeOf((*MockIColoursDomain)(nil).GetColourRoleName), arg0, arg1)
}

// GetHistory mocks base method.
func (m *MockIColoursDomain) GetHistory(arg0 context.Context, arg1 IDomainMember) (*History, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistory", arg0, arg1)
	ret0, _ := ret[0].(*History)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHistory indicates an expected call of GetHistory.
func (mr *MockIColoursDomainMockRecorder) GetHistory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistory", reflect.TypeOf((*MockIColoursDomain)(nil).GetHistory), arg0, arg1)
}

// GetLastFrozen mocks base method.
func (m *MockIColoursDomain) GetLastFrozen(arg0 context.Context, arg1 IDomainMember) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastFrozen", arg0, arg1)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastFrozen indicates an expected call of GetLastFrozen.
func (mr *MockIColoursDomainMockRecorder) GetLastFrozen(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastFrozen", reflect.TypeOf((*MockIColoursDomain)(nil).GetLastFrozen), arg0, arg1)
}

// GetLastMutate mocks base method.
func (m *MockIColoursDomain) GetLastMutate(arg0 context.Context, arg1 IDomainMember) (time.Time, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastMutate", arg0, arg1)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLastMutate indicates an expected call of GetLastMutate.
func (mr *MockIColoursDomainMockRecorder) GetLastMutate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastMutate", reflect.TypeOf((*MockIColoursDomain)(nil).GetLastMutate), arg0, arg1)
}

// GetLastReroll mocks base method.
func (m *MockIColoursDomain) GetLastReroll(arg0 context.Context, arg1 IDomainMember) (time.Time, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastReroll", arg0, arg1)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLastReroll indicates an expected call of GetLastReroll.
func (mr *MockIColoursDomainMockRecorder) GetLastReroll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastReroll", reflect.TypeOf((*MockIColoursDomain)(nil).GetLastReroll), arg0, arg1)
}

// GetRerollCooldownEndTime mocks base method.
func (m *MockIColoursDomain) GetRerollCooldownEndTime(arg0 context.Context, arg1 IDomainMember) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRerollCooldownEndTime", arg0, arg1)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRerollCooldownEndTime indicates an expected call of GetRerollCooldownEndTime.
func (mr *MockIColoursDomainMockRecorder) GetRerollCooldownEndTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRerollCooldownEndTime", reflect.TypeOf((*MockIColoursDomain)(nil).GetRerollCooldownEndTime), arg0, arg1)
}

// HasColourRole mocks base method.
func (m *MockIColoursDomain) HasColourRole(arg0 context.Context, arg1 IDomainMember) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasColourRole", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasColourRole indicates an expected call of HasColourRole.
func (mr *MockIColoursDomainMockRecorder) HasColourRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasColourRole", reflect.TypeOf((*MockIColoursDomain)(nil).HasColourRole), arg0, arg1)
}

// Mutate mocks base method.
func (m *MockIColoursDomain) Mutate(arg0 context.Context, arg1 IDomainSession, arg2 IDomainMember) (*Colour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mutate", arg0, arg1, arg2)
	ret0, _ := ret[0].(*Colour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Mutate indicates an expected call of Mutate.
func (mr *MockIColoursDomainMockRecorder) Mutate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mutate", reflect.TypeOf((*MockIColoursDomain)(nil).Mutate), arg0, arg1, arg2)
}

// Reroll mocks base method.
func (m *MockIColoursDomain) Reroll(arg0 context.Context, arg1 IDomainSession, arg2 IDomainMember) (*Colour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reroll", arg0, arg1, arg2)
	ret0, _ := ret[0].(*Colour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reroll indicates an expected call of Reroll.
func (mr *MockIColoursDomainMockRecorder) Reroll(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reroll", reflect.TypeOf((*MockIColoursDomain)(nil).Reroll), arg0, arg1, arg2)
}

// SetRoleHeight mocks base method.
func (m *MockIColoursDomain) SetRoleHeight(arg0 context.Context, arg1 IDomainSession, arg2 IDomainGuild, arg3 string, arg4 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRoleHeight", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRoleHeight indicates an expected call of SetRoleHeight.
func (mr *MockIColoursDomainMockRecorder) SetRoleHeight(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRoleHeight", reflect.TypeOf((*MockIColoursDomain)(nil).SetRoleHeight), arg0, arg1, arg2, arg3, arg4)
}

// Unfreeze mocks base method.
func (m *MockIColoursDomain) Unfreeze(arg0 context.Context, arg1 IDomainMember) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unfreeze", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unfreeze indicates an expected call of Unfreeze.
func (mr *MockIColoursDomainMockRecorder) Unfreeze(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unfreeze", reflect.TypeOf((*MockIColoursDomain)(nil).Unfreeze), arg0, arg1)
}

// MockIDomainSession is a mock of IDomainSession interface.
type MockIDomainSession struct {
	ctrl     *gomock.Controller
	recorder *MockIDomainSessionMockRecorder
}

// MockIDomainSessionMockRecorder is the mock recorder for MockIDomainSession.
type MockIDomainSessionMockRecorder struct {
	mock *MockIDomainSession
}

// NewMockIDomainSession creates a new mock instance.
func NewMockIDomainSession(ctrl *gomock.Controller) *MockIDomainSession {
	mock := &MockIDomainSession{ctrl: ctrl}
	mock.recorder = &MockIDomainSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDomainSession) EXPECT() *MockIDomainSessionMockRecorder {
	return m.recorder
}

// GuildMember mocks base method.
func (m *MockIDomainSession) GuildMember(ctx context.Context, guildID, userID string) (IDomainMember, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuildMember", ctx, guildID, userID)
	ret0, _ := ret[0].(IDomainMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GuildMember indicates an expected call of GuildMember.
func (mr *MockIDomainSessionMockRecorder) GuildMember(ctx, guildID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuildMember", reflect.TypeOf((*MockIDomainSession)(nil).GuildMember), ctx, guildID, userID)
}

// GuildMemberRoleAdd mocks base method.
func (m *MockIDomainSession) GuildMemberRoleAdd(ctx context.Context, guildID, userID, roleID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuildMemberRoleAdd", ctx, guildID, userID, roleID)
	ret0, _ := ret[0].(error)
	return ret0
}

// GuildMemberRoleAdd indicates an expected call of GuildMemberRoleAdd.
func (mr *MockIDomainSessionMockRecorder) GuildMemberRoleAdd(ctx, guildID, userID, roleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuildMemberRoleAdd", reflect.TypeOf((*MockIDomainSession)(nil).GuildMemberRoleAdd), ctx, guildID, userID, roleID)
}

// GuildRole mocks base method.
func (m *MockIDomainSession) GuildRole(ctx context.Context, guildID, roleID string) (IDomainRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuildRole", ctx, guildID, roleID)
	ret0, _ := ret[0].(IDomainRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GuildRole indicates an expected call of GuildRole.
func (mr *MockIDomainSessionMockRecorder) GuildRole(ctx, guildID, roleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuildRole", reflect.TypeOf((*MockIDomainSession)(nil).GuildRole), ctx, guildID, roleID)
}

// GuildRoleCreate mocks base method.
func (m *MockIDomainSession) GuildRoleCreate(ctx context.Context, guildID, name string, colour int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuildRoleCreate", ctx, guildID, name, colour)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GuildRoleCreate indicates an expected call of GuildRoleCreate.
func (mr *MockIDomainSessionMockRecorder) GuildRoleCreate(ctx, guildID, name, colour interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuildRoleCreate", reflect.TypeOf((*MockIDomainSession)(nil).GuildRoleCreate), ctx, guildID, name, colour)
}

// GuildRoleEdit mocks base method.
func (m *MockIDomainSession) GuildRoleEdit(ctx context.Context, guildID, roleID, name string, colour int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuildRoleEdit", ctx, guildID, roleID, name, colour)
	ret0, _ := ret[0].(error)
	return ret0
}

// GuildRoleEdit indicates an expected call of GuildRoleEdit.
func (mr *MockIDomainSessionMockRecorder) GuildRoleEdit(ctx, guildID, roleID, name, colour interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuildRoleEdit", reflect.TypeOf((*MockIDomainSession)(nil).GuildRoleEdit), ctx, guildID, roleID, name, colour)
}

// GuildRoleReorder mocks base method.
func (m *MockIDomainSession) GuildRoleReorder(ctx context.Context, guildID string, roles []IDomainRole) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuildRoleReorder", ctx, guildID, roles)
	ret0, _ := ret[0].(error)
	return ret0
}

// GuildRoleReorder indicates an expected call of GuildRoleReorder.
func (mr *MockIDomainSessionMockRecorder) GuildRoleReorder(ctx, guildID, roles interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuildRoleReorder", reflect.TypeOf((*MockIDomainSession)(nil).GuildRoleReorder), ctx, guildID, roles)
}

// GuildRoles mocks base method.
func (m *MockIDomainSession) GuildRoles(ctx context.Context, guildID string) ([]IDomainRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuildRoles", ctx, guildID)
	ret0, _ := ret[0].([]IDomainRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GuildRoles indicates an expected call of GuildRoles.
func (mr *MockIDomainSessionMockRecorder) GuildRoles(ctx, guildID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuildRoles", reflect.TypeOf((*MockIDomainSession)(nil).GuildRoles), ctx, guildID)
}

// MockIDomainGuild is a mock of IDomainGuild interface.
type MockIDomainGuild struct {
	ctrl     *gomock.Controller
	recorder *MockIDomainGuildMockRecorder
}

// MockIDomainGuildMockRecorder is the mock recorder for MockIDomainGuild.
type MockIDomainGuildMockRecorder struct {
	mock *MockIDomainGuild
}

// NewMockIDomainGuild creates a new mock instance.
func NewMockIDomainGuild(ctrl *gomock.Controller) *MockIDomainGuild {
	mock := &MockIDomainGuild{ctrl: ctrl}
	mock.recorder = &MockIDomainGuildMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDomainGuild) EXPECT() *MockIDomainGuildMockRecorder {
	return m.recorder
}

// ID mocks base method.
func (m *MockIDomainGuild) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockIDomainGuildMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockIDomainGuild)(nil).ID))
}

// MockIDomainMember is a mock of IDomainMember interface.
type MockIDomainMember struct {
	ctrl     *gomock.Controller
	recorder *MockIDomainMemberMockRecorder
}

// MockIDomainMemberMockRecorder is the mock recorder for MockIDomainMember.
type MockIDomainMemberMockRecorder struct {
	mock *MockIDomainMember
}

// NewMockIDomainMember creates a new mock instance.
func NewMockIDomainMember(ctrl *gomock.Controller) *MockIDomainMember {
	mock := &MockIDomainMember{ctrl: ctrl}
	mock.recorder = &MockIDomainMemberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDomainMember) EXPECT() *MockIDomainMemberMockRecorder {
	return m.recorder
}

// CacheKey mocks base method.
func (m *MockIDomainMember) CacheKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CacheKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// CacheKey indicates an expected call of CacheKey.
func (mr *MockIDomainMemberMockRecorder) CacheKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CacheKey", reflect.TypeOf((*MockIDomainMember)(nil).CacheKey))
}

// Guild mocks base method.
func (m *MockIDomainMember) Guild() IDomainGuild {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Guild")
	ret0, _ := ret[0].(IDomainGuild)
	return ret0
}

// Guild indicates an expected call of Guild.
func (mr *MockIDomainMemberMockRecorder) Guild() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Guild", reflect.TypeOf((*MockIDomainMember)(nil).Guild))
}

// Nick mocks base method.
func (m *MockIDomainMember) Nick() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Nick")
	ret0, _ := ret[0].(string)
	return ret0
}

// Nick indicates an expected call of Nick.
func (mr *MockIDomainMemberMockRecorder) Nick() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nick", reflect.TypeOf((*MockIDomainMember)(nil).Nick))
}

// Roles mocks base method.
func (m *MockIDomainMember) Roles() []IDomainRole {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Roles")
	ret0, _ := ret[0].([]IDomainRole)
	return ret0
}

// Roles indicates an expected call of Roles.
func (mr *MockIDomainMemberMockRecorder) Roles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Roles", reflect.TypeOf((*MockIDomainMember)(nil).Roles))
}

// UserID mocks base method.
func (m *MockIDomainMember) UserID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserID")
	ret0, _ := ret[0].(string)
	return ret0
}

// UserID indicates an expected call of UserID.
func (mr *MockIDomainMemberMockRecorder) UserID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserID", reflect.TypeOf((*MockIDomainMember)(nil).UserID))
}

// Username mocks base method.
func (m *MockIDomainMember) Username() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Username")
	ret0, _ := ret[0].(string)
	return ret0
}

// Username indicates an expected call of Username.
func (mr *MockIDomainMemberMockRecorder) Username() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Username", reflect.TypeOf((*MockIDomainMember)(nil).Username))
}

// MockIDomainRole is a mock of IDomainRole interface.
type MockIDomainRole struct {
	ctrl     *gomock.Controller
	recorder *MockIDomainRoleMockRecorder
}

// MockIDomainRoleMockRecorder is the mock recorder for MockIDomainRole.
type MockIDomainRoleMockRecorder struct {
	mock *MockIDomainRole
}

// NewMockIDomainRole creates a new mock instance.
func NewMockIDomainRole(ctrl *gomock.Controller) *MockIDomainRole {
	mock := &MockIDomainRole{ctrl: ctrl}
	mock.recorder = &MockIDomainRoleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDomainRole) EXPECT() *MockIDomainRoleMockRecorder {
	return m.recorder
}

// CacheKey mocks base method.
func (m *MockIDomainRole) CacheKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CacheKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// CacheKey indicates an expected call of CacheKey.
func (mr *MockIDomainRoleMockRecorder) CacheKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CacheKey", reflect.TypeOf((*MockIDomainRole)(nil).CacheKey))
}

// Colour mocks base method.
func (m *MockIDomainRole) Colour() *Colour {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Colour")
	ret0, _ := ret[0].(*Colour)
	return ret0
}

// Colour indicates an expected call of Colour.
func (mr *MockIDomainRoleMockRecorder) Colour() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Colour", reflect.TypeOf((*MockIDomainRole)(nil).Colour))
}

// ID mocks base method.
func (m *MockIDomainRole) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockIDomainRoleMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockIDomainRole)(nil).ID))
}

// Name mocks base method.
func (m *MockIDomainRole) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockIDomainRoleMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockIDomainRole)(nil).Name))
}

// MockIDomainRepository is a mock of IDomainRepository interface.
type MockIDomainRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIDomainRepositoryMockRecorder
}

// MockIDomainRepositoryMockRecorder is the mock recorder for MockIDomainRepository.
type MockIDomainRepositoryMockRecorder struct {
	mock *MockIDomainRepository
}

// NewMockIDomainRepository creates a new mock instance.
func NewMockIDomainRepository(ctrl *gomock.Controller) *MockIDomainRepository {
	mock := &MockIDomainRepository{ctrl: ctrl}
	mock.recorder = &MockIDomainRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDomainRepository) EXPECT() *MockIDomainRepositoryMockRecorder {
	return m.recorder
}

// FetchUserHistory mocks base method.
func (m *MockIDomainRepository) FetchUserHistory(arg0 context.Context, arg1 IDomainMember, arg2 time.Time) ([]*ColoursLogRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchUserHistory", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*ColoursLogRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchUserHistory indicates an expected call of FetchUserHistory.
func (mr *MockIDomainRepositoryMockRecorder) FetchUserHistory(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUserHistory", reflect.TypeOf((*MockIDomainRepository)(nil).FetchUserHistory), arg0, arg1, arg2)
}

// FetchUserState mocks base method.
func (m *MockIDomainRepository) FetchUserState(arg0 context.Context, arg1 IDomainMember, arg2 Reason) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchUserState", arg0, arg1, arg2)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchUserState indicates an expected call of FetchUserState.
func (mr *MockIDomainRepositoryMockRecorder) FetchUserState(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUserState", reflect.TypeOf((*MockIDomainRepository)(nil).FetchUserState), arg0, arg1, arg2)
}

// UpdateFreeze mocks base method.
func (m *MockIDomainRepository) UpdateFreeze(arg0 context.Context, arg1 IDomainMember) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFreeze", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFreeze indicates an expected call of UpdateFreeze.
func (mr *MockIDomainRepositoryMockRecorder) UpdateFreeze(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFreeze", reflect.TypeOf((*MockIDomainRepository)(nil).UpdateFreeze), arg0, arg1)
}

// UpdateMutate mocks base method.
func (m *MockIDomainRepository) UpdateMutate(arg0 context.Context, arg1 IDomainMember, arg2 *Colour) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMutate", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMutate indicates an expected call of UpdateMutate.
func (mr *MockIDomainRepositoryMockRecorder) UpdateMutate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMutate", reflect.TypeOf((*MockIDomainRepository)(nil).UpdateMutate), arg0, arg1, arg2)
}

// UpdateReroll mocks base method.
func (m *MockIDomainRepository) UpdateReroll(arg0 context.Context, arg1 IDomainMember, arg2 *Colour) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateReroll", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateReroll indicates an expected call of UpdateReroll.
func (mr *MockIDomainRepositoryMockRecorder) UpdateReroll(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateReroll", reflect.TypeOf((*MockIDomainRepository)(nil).UpdateReroll), arg0, arg1, arg2)
}

// UpdateRerollPenalty mocks base method.
func (m *MockIDomainRepository) UpdateRerollPenalty(arg0 context.Context, arg1 IDomainMember, arg2 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRerollPenalty", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRerollPenalty indicates an expected call of UpdateRerollPenalty.
func (mr *MockIDomainRepositoryMockRecorder) UpdateRerollPenalty(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRerollPenalty", reflect.TypeOf((*MockIDomainRepository)(nil).UpdateRerollPenalty), arg0, arg1, arg2)
}

// UpdateUnfreeze mocks base method.
func (m *MockIDomainRepository) UpdateUnfreeze(arg0 context.Context, arg1 IDomainMember) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUnfreeze", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUnfreeze indicates an expected call of UpdateUnfreeze.
func (mr *MockIDomainRepositoryMockRecorder) UpdateUnfreeze(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUnfreeze", reflect.TypeOf((*MockIDomainRepository)(nil).UpdateUnfreeze), arg0, arg1)
}
