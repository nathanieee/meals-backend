// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	models "project-skbackend/internal/models"

	mock "github.com/stretchr/testify/mock"

	requests "project-skbackend/internal/controllers/requests"

	responses "project-skbackend/internal/controllers/responses"

	utpagination "project-skbackend/packages/utils/utpagination"

	uuid "github.com/google/uuid"
)

// IPartnerService is an autogenerated mock type for the IPartnerService type
type IPartnerService struct {
	mock.Mock
}

// Create provides a mock function with given fields: req
func (_m *IPartnerService) Create(req requests.CreatePartner) (*responses.Partner, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *responses.Partner
	var r1 error
	if rf, ok := ret.Get(0).(func(requests.CreatePartner) (*responses.Partner, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(requests.CreatePartner) *responses.Partner); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*responses.Partner)
		}
	}

	if rf, ok := ret.Get(1).(func(requests.CreatePartner) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *IPartnerService) Delete(id uuid.UUID) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: preq
func (_m *IPartnerService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	ret := _m.Called(preq)

	if len(ret) == 0 {
		panic("no return value specified for FindAll")
	}

	var r0 *utpagination.Pagination
	var r1 error
	if rf, ok := ret.Get(0).(func(utpagination.Pagination) (*utpagination.Pagination, error)); ok {
		return rf(preq)
	}
	if rf, ok := ret.Get(0).(func(utpagination.Pagination) *utpagination.Pagination); ok {
		r0 = rf(preq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utpagination.Pagination)
		}
	}

	if rf, ok := ret.Get(1).(func(utpagination.Pagination) error); ok {
		r1 = rf(preq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByID provides a mock function with given fields: id
func (_m *IPartnerService) FindByID(id uuid.UUID) (*responses.Partner, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *responses.Partner
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*responses.Partner, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *responses.Partner); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*responses.Partner)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Read provides a mock function with given fields:
func (_m *IPartnerService) Read() ([]*models.Partner, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 []*models.Partner
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*models.Partner, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*models.Partner); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Partner)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, req
func (_m *IPartnerService) Update(id uuid.UUID, req requests.UpdatePartner) (*responses.Partner, error) {
	ret := _m.Called(id, req)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *responses.Partner
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, requests.UpdatePartner) (*responses.Partner, error)); ok {
		return rf(id, req)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, requests.UpdatePartner) *responses.Partner); ok {
		r0 = rf(id, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*responses.Partner)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, requests.UpdatePartner) error); ok {
		r1 = rf(id, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIPartnerService creates a new instance of IPartnerService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIPartnerService(t interface {
	mock.TestingT
	Cleanup(func())
}) *IPartnerService {
	mock := &IPartnerService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}