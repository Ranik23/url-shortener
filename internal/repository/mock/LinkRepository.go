// Code generated by mockery v2.53.3. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// LinkRepository is an autogenerated mock type for the LinkRepository type
type LinkRepository struct {
	mock.Mock
}

// CreateLink provides a mock function with given fields: ctx, default_link, shortened_link
func (_m *LinkRepository) CreateLink(ctx context.Context, default_link string, shortened_link string) error {
	ret := _m.Called(ctx, default_link, shortened_link)

	if len(ret) == 0 {
		panic("no return value specified for CreateLink")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, default_link, shortened_link)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteLink provides a mock function with given fields: ctx, default_link
func (_m *LinkRepository) DeleteLink(ctx context.Context, default_link string) error {
	ret := _m.Called(ctx, default_link)

	if len(ret) == 0 {
		panic("no return value specified for DeleteLink")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, default_link)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDefaultLink provides a mock function with given fields: ctx, shortened_link
func (_m *LinkRepository) GetDefaultLink(ctx context.Context, shortened_link string) (string, error) {
	ret := _m.Called(ctx, shortened_link)

	if len(ret) == 0 {
		panic("no return value specified for GetDefaultLink")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, shortened_link)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, shortened_link)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, shortened_link)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShortenedLink provides a mock function with given fields: ctx, default_link
func (_m *LinkRepository) GetShortenedLink(ctx context.Context, default_link string) (string, error) {
	ret := _m.Called(ctx, default_link)

	if len(ret) == 0 {
		panic("no return value specified for GetShortenedLink")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, default_link)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, default_link)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, default_link)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewLinkRepository creates a new instance of LinkRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLinkRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *LinkRepository {
	mock := &LinkRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
