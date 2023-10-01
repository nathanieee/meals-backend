package repositories

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"

	"github.com/google/uuid"
)

type (
	IUserRepo interface {
		FindAll(p models.Pagination) (*models.Pagination, error)
		Create(user *models.User) (*models.User, error)
		Update(user models.User, uid uuid.UUID) (*models.User, error)
		FindByID(uid uuid.UUID) (*responses.UserResponse, error)
		FindByEmail(email string) (*responses.UserResponse, error)
		Delete(user models.User) error
	}

	IMemberRepo interface {
		Create(m *models.Member) (*models.Member, error)
		Update(m models.Member, mid uuid.UUID) (*models.Member, error)
		FindAll(p models.Pagination) (*models.Pagination, error)
		FindByID(mid uuid.UUID) (*responses.MemberResponse, error)
		Delete(m models.Member) error
	}

	IAdminRepo interface {
		Create(a *models.Admin) (*models.Admin, error)
		Update(a models.Admin, aid uuid.UUID) (*models.Admin, error)
		FindAll(p models.Pagination) (*models.Pagination, error)
		FindByID(aid uuid.UUID) (*responses.AdminResponse, error)
		Delete(a models.Admin) error
	}

	ICaregiverRepo interface {
		Create(cg *models.Caregiver) (*models.Caregiver, error)
		Update(cg models.Caregiver, cgid uuid.UUID) (*models.Caregiver, error)
		FindAll(p models.Pagination) (*models.Pagination, error)
		FindByID(cgid uuid.UUID) (*responses.CaregiverResponse, error)
		Delete(cg models.Caregiver) error
	}

	IMealRepo interface{}
)
