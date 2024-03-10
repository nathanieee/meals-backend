// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	models "project-skbackend/internal/models"

	mock "github.com/stretchr/testify/mock"

	utpagination "project-skbackend/packages/utils/utpagination"

	uuid "github.com/google/uuid"
)

// ICaregiverRepository is an autogenerated mock type for the ICaregiverRepository type
type ICaregiverRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: cg
func (_m *ICaregiverRepository) Create(cg models.Caregiver) (*models.Caregiver, error) {
	ret := _m.Called(cg)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.Caregiver
	var r1 error
	if rf, ok := ret.Get(0).(func(models.Caregiver) (*models.Caregiver, error)); ok {
		return rf(cg)
	}
	if rf, ok := ret.Get(0).(func(models.Caregiver) *models.Caregiver); ok {
		r0 = rf(cg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Caregiver)
		}
	}

	if rf, ok := ret.Get(1).(func(models.Caregiver) error); ok {
		r1 = rf(cg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: cg
func (_m *ICaregiverRepository) Delete(cg models.Caregiver) error {
	ret := _m.Called(cg)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Caregiver) error); ok {
		r0 = rf(cg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: p
func (_m *ICaregiverRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for FindAll")
	}

	var r0 *utpagination.Pagination
	var r1 error
	if rf, ok := ret.Get(0).(func(utpagination.Pagination) (*utpagination.Pagination, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func(utpagination.Pagination) *utpagination.Pagination); ok {
		r0 = rf(p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utpagination.Pagination)
		}
	}

	if rf, ok := ret.Get(1).(func(utpagination.Pagination) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByEmail provides a mock function with given fields: email
func (_m *ICaregiverRepository) FindByEmail(email string) (*models.Caregiver, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for FindByEmail")
	}

	var r0 *models.Caregiver
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.Caregiver, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *models.Caregiver); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Caregiver)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByID provides a mock function with given fields: cgid
func (_m *ICaregiverRepository) FindByID(cgid uuid.UUID) (*models.Caregiver, error) {
	ret := _m.Called(cgid)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *models.Caregiver
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.Caregiver, error)); ok {
		return rf(cgid)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Caregiver); ok {
		r0 = rf(cgid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Caregiver)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(cgid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Read provides a mock function with given fields:
func (_m *ICaregiverRepository) Read() ([]*models.Caregiver, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 []*models.Caregiver
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*models.Caregiver, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*models.Caregiver); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Caregiver)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: cg
func (_m *ICaregiverRepository) Update(cg models.Caregiver) (*models.Caregiver, error) {
	ret := _m.Called(cg)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *models.Caregiver
	var r1 error
	if rf, ok := ret.Get(0).(func(models.Caregiver) (*models.Caregiver, error)); ok {
		return rf(cg)
	}
	if rf, ok := ret.Get(0).(func(models.Caregiver) *models.Caregiver); ok {
		r0 = rf(cg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Caregiver)
		}
	}

	if rf, ok := ret.Get(1).(func(models.Caregiver) error); ok {
		r1 = rf(cg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewICaregiverRepository creates a new instance of ICaregiverRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICaregiverRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICaregiverRepository {
	mock := &ICaregiverRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}