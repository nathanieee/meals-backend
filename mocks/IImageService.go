// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	consttypes "project-skbackend/packages/consttypes"

	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"

	multipart "mime/multipart"
)

// IImageService is an autogenerated mock type for the IImageService type
type IImageService struct {
	mock.Mock
}

// Upload provides a mock function with given fields: fileheader, imgtype, ctx
func (_m *IImageService) Upload(fileheader *multipart.FileHeader, imgtype consttypes.ImageType, ctx *gin.Context) error {
	ret := _m.Called(fileheader, imgtype, ctx)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*multipart.FileHeader, consttypes.ImageType, *gin.Context) error); ok {
		r0 = rf(fileheader, imgtype, ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIImageService creates a new instance of IImageService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIImageService(t interface {
	mock.TestingT
	Cleanup(func())
}) *IImageService {
	mock := &IImageService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}