// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// AddPoint provides a mock function with given fields: z
func (_m *Storage) AddPoint(z uint64) {
	_m.Called(z)
}

// Len provides a mock function with given fields:
func (_m *Storage) Len() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// RLock provides a mock function with given fields:
func (_m *Storage) RLock() {
	_m.Called()
}

// RUnlock provides a mock function with given fields:
func (_m *Storage) RUnlock() {
	_m.Called()
}

// ReadPoint provides a mock function with given fields: z
func (_m *Storage) ReadPoint(z uint64) (bool, bool) {
	ret := _m.Called(z)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint64) bool); ok {
		r0 = rf(z)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(uint64) bool); ok {
		r1 = rf(z)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}
