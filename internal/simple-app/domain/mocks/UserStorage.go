// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UserStorage is an autogenerated mock type for the UserStorage type
type UserStorage struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *UserStorage) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Read provides a mock function with given fields: ctx, id
func (_m *UserStorage) Read(ctx context.Context, id string) (string, error) {
	ret := _m.Called(ctx, id)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, id, name
func (_m *UserStorage) Store(ctx context.Context, id string, name string) error {
	ret := _m.Called(ctx, id, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, id, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserStorage creates a new instance of UserStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserStorage {
	mock := &UserStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
