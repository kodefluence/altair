// Code generated by MockGen. DO NOT EDIT.
// Source: ./http.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/kodefluence/altair/plugin/oauth/entity"
	jsonapi "github.com/kodefluence/monorepo/jsonapi"
	kontext "github.com/kodefluence/monorepo/kontext"
)

// MockApplicationManager is a mock of ApplicationManager interface.
type MockApplicationManager struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationManagerMockRecorder
}

// MockApplicationManagerMockRecorder is the mock recorder for MockApplicationManager.
type MockApplicationManagerMockRecorder struct {
	mock *MockApplicationManager
}

// NewMockApplicationManager creates a new mock instance.
func NewMockApplicationManager(ctrl *gomock.Controller) *MockApplicationManager {
	mock := &MockApplicationManager{ctrl: ctrl}
	mock.recorder = &MockApplicationManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplicationManager) EXPECT() *MockApplicationManagerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockApplicationManager) Create(ktx kontext.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, jsonapi.Errors) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ktx, e)
	ret0, _ := ret[0].(entity.OauthApplicationJSON)
	ret1, _ := ret[1].(jsonapi.Errors)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockApplicationManagerMockRecorder) Create(ktx, e interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockApplicationManager)(nil).Create), ktx, e)
}

// List mocks base method.
func (m *MockApplicationManager) List(ktx kontext.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, jsonapi.Errors) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ktx, offset, limit)
	ret0, _ := ret[0].([]entity.OauthApplicationJSON)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(jsonapi.Errors)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockApplicationManagerMockRecorder) List(ktx, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockApplicationManager)(nil).List), ktx, offset, limit)
}

// One mocks base method.
func (m *MockApplicationManager) One(ktx kontext.Context, ID int) (entity.OauthApplicationJSON, jsonapi.Errors) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "One", ktx, ID)
	ret0, _ := ret[0].(entity.OauthApplicationJSON)
	ret1, _ := ret[1].(jsonapi.Errors)
	return ret0, ret1
}

// One indicates an expected call of One.
func (mr *MockApplicationManagerMockRecorder) One(ktx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "One", reflect.TypeOf((*MockApplicationManager)(nil).One), ktx, ID)
}

// Update mocks base method.
func (m *MockApplicationManager) Update(ktx kontext.Context, ID int, e entity.OauthApplicationUpdateJSON) (entity.OauthApplicationJSON, jsonapi.Errors) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ktx, ID, e)
	ret0, _ := ret[0].(entity.OauthApplicationJSON)
	ret1, _ := ret[1].(jsonapi.Errors)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockApplicationManagerMockRecorder) Update(ktx, ID, e interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockApplicationManager)(nil).Update), ktx, ID, e)
}
