// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	aws "github.com/applike/gosoline/pkg/cloud/aws"

	mock "github.com/stretchr/testify/mock"
)

// Executor is an autogenerated mock type for the Executor type
type Executor struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, f
func (_m *Executor) Execute(ctx context.Context, f aws.RequestFunction) (interface{}, error) {
	ret := _m.Called(ctx, f)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, aws.RequestFunction) interface{}); ok {
		r0 = rf(ctx, f)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, aws.RequestFunction) error); ok {
		r1 = rf(ctx, f)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}