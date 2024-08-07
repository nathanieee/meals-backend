// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	requests "project-skbackend/internal/controllers/requests"
	responses "project-skbackend/internal/controllers/responses"

	mock "github.com/stretchr/testify/mock"

	utpagination "project-skbackend/packages/utils/utpagination"

	uuid "github.com/google/uuid"
)

// IMemberService is an autogenerated mock type for the IMemberService type
type IMemberService struct {
	mock.Mock
}

// Create provides a mock function with given fields: req
func (_m *IMemberService) Create(req requests.CreateMember) (*responses.Member, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *responses.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(requests.CreateMember) (*responses.Member, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(requests.CreateMember) *responses.Member); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*responses.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(requests.CreateMember) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *IMemberService) Delete(id uuid.UUID) error {
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
func (_m *IMemberService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
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

// GetByCaregiverID provides a mock function with given fields: cgid
func (_m *IMemberService) GetByCaregiverID(cgid uuid.UUID) (*responses.Member, error) {
	ret := _m.Called(cgid)

	if len(ret) == 0 {
		panic("no return value specified for GetByCaregiverID")
	}

	var r0 *responses.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*responses.Member, error)); ok {
		return rf(cgid)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *responses.Member); ok {
		r0 = rf(cgid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*responses.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(cgid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *IMemberService) GetByID(id uuid.UUID) (*responses.Member, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *responses.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*responses.Member, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *responses.Member); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*responses.Member)
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
func (_m *IMemberService) Read() ([]*responses.Member, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 []*responses.Member
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*responses.Member, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*responses.Member); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*responses.Member)
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
func (_m *IMemberService) Update(id uuid.UUID, req requests.UpdateMember) (*responses.Member, error) {
	ret := _m.Called(id, req)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *responses.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, requests.UpdateMember) (*responses.Member, error)); ok {
		return rf(id, req)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, requests.UpdateMember) *responses.Member); ok {
		r0 = rf(id, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*responses.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, requests.UpdateMember) error); ok {
		r1 = rf(id, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIMemberService creates a new instance of IMemberService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIMemberService(t interface {
	mock.TestingT
	Cleanup(func())
}) *IMemberService {
	mock := &IMemberService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
