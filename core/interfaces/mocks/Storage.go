// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	interfaces "application/core/interfaces"

	mock "github.com/stretchr/testify/mock"

	model "application/core/model"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// DeleteConfig provides a mock function with given fields: id
func (_m *Storage) DeleteConfig(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteExample provides a mock function with given fields: orgID, appID, id
func (_m *Storage) DeleteExample(orgID string, appID string, id string) error {
	ret := _m.Called(orgID, appID, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(orgID, appID, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindExample provides a mock function with given fields: orgID, appID, id
func (_m *Storage) FindExample(orgID string, appID string, id string) (*model.Example, error) {
	ret := _m.Called(orgID, appID, id)

	var r0 *model.Example
	if rf, ok := ret.Get(0).(func(string, string, string) *model.Example); ok {
		r0 = rf(orgID, appID, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Example)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(orgID, appID, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConfig provides a mock function with given fields: id
func (_m *Storage) GetConfig(id string) (*model.Config, error) {
	ret := _m.Called(id)

	var r0 *model.Config
	if rf, ok := ret.Get(0).(func(string) *model.Config); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Config)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertExample provides a mock function with given fields: example
func (_m *Storage) InsertExample(example model.Example) error {
	ret := _m.Called(example)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.Example) error); ok {
		r0 = rf(example)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PerformTransaction provides a mock function with given fields: _a0
func (_m *Storage) PerformTransaction(_a0 func(interfaces.Storage) error) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(func(interfaces.Storage) error) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterStorageListener provides a mock function with given fields: listener
func (_m *Storage) RegisterStorageListener(listener interfaces.StorageListener) {
	_m.Called(listener)
}

// SaveConfig provides a mock function with given fields: configs
func (_m *Storage) SaveConfig(configs model.Config) error {
	ret := _m.Called(configs)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.Config) error); ok {
		r0 = rf(configs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateExample provides a mock function with given fields: example
func (_m *Storage) UpdateExample(example model.Example) error {
	ret := _m.Called(example)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.Example) error); ok {
		r0 = rf(example)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorage(t mockConstructorTestingTNewStorage) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}