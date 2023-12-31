// Code generated by MockGen. DO NOT EDIT.
// Source: iface/person_repo.go

// Package mock_iface is a generated GoMock package.
package mock_iface

import (
	model "fio-service/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPersonRepo is a mock of PersonRepo interface.
type MockPersonRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPersonRepoMockRecorder
}

// MockPersonRepoMockRecorder is the mock recorder for MockPersonRepo.
type MockPersonRepoMockRecorder struct {
	mock *MockPersonRepo
}

// NewMockPersonRepo creates a new mock instance.
func NewMockPersonRepo(ctrl *gomock.Controller) *MockPersonRepo {
	mock := &MockPersonRepo{ctrl: ctrl}
	mock.recorder = &MockPersonRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersonRepo) EXPECT() *MockPersonRepoMockRecorder {
	return m.recorder
}

// AddPerson mocks base method.
func (m *MockPersonRepo) AddPerson(arg0 *model.Person) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPerson", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPerson indicates an expected call of AddPerson.
func (mr *MockPersonRepoMockRecorder) AddPerson(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPerson", reflect.TypeOf((*MockPersonRepo)(nil).AddPerson), arg0)
}

// DeletePersonById mocks base method.
func (m *MockPersonRepo) DeletePersonById(id int) (*model.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePersonById", id)
	ret0, _ := ret[0].(*model.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeletePersonById indicates an expected call of DeletePersonById.
func (mr *MockPersonRepoMockRecorder) DeletePersonById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePersonById", reflect.TypeOf((*MockPersonRepo)(nil).DeletePersonById), id)
}

// GetPeopleByFilters mocks base method.
func (m *MockPersonRepo) GetPeopleByFilters(arg0 *model.Filters) ([]*model.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeopleByFilters", arg0)
	ret0, _ := ret[0].([]*model.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPeopleByFilters indicates an expected call of GetPeopleByFilters.
func (mr *MockPersonRepoMockRecorder) GetPeopleByFilters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeopleByFilters", reflect.TypeOf((*MockPersonRepo)(nil).GetPeopleByFilters), arg0)
}

// GetPersonById mocks base method.
func (m *MockPersonRepo) GetPersonById(id int) (*model.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPersonById", id)
	ret0, _ := ret[0].(*model.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPersonById indicates an expected call of GetPersonById.
func (mr *MockPersonRepoMockRecorder) GetPersonById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPersonById", reflect.TypeOf((*MockPersonRepo)(nil).GetPersonById), id)
}

// UpdatePerson mocks base method.
func (m *MockPersonRepo) UpdatePerson(arg0 *model.UpdatePerson) (*model.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePerson", arg0)
	ret0, _ := ret[0].(*model.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePerson indicates an expected call of UpdatePerson.
func (mr *MockPersonRepoMockRecorder) UpdatePerson(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePerson", reflect.TypeOf((*MockPersonRepo)(nil).UpdatePerson), arg0)
}
