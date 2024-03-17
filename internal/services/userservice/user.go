package userservice

import (
	"encoding/json"
	"fmt"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/base"
	"project-skbackend/internal/repositories/adminrepo"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/organizationrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utpagination"
	"project-skbackend/packages/utils/utresponse"

	"github.com/google/uuid"
)

type (
	UserService struct {
		ruser userrepo.IUserRepository
		radmn adminrepo.IAdminRepository
		rcare caregiverrepo.ICaregiverRepository
		rmemb memberrepo.IMemberRepository
		rorga organizationrepo.IOrganizationRepository
		rpart partnerrepo.IPartnerRepository
		// rpatr patronrepo.IPatronRepository
	}

	IUserService interface {
		Create(req requests.CreateUser) (*responses.User, error)
		FindByID(uid uuid.UUID) (*responses.User, error)
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		Delete(uid uuid.UUID) error
		Update(req requests.UpdateUser, uid uuid.UUID) (*responses.User, error)
	}
)

func NewUserService(
	ruser userrepo.IUserRepository,
	radmn adminrepo.IAdminRepository,
	rcare caregiverrepo.ICaregiverRepository,
	rmemb memberrepo.IMemberRepository,
	rorga organizationrepo.IOrganizationRepository,
	rpart partnerrepo.IPartnerRepository,
	// rpatr patronrepo.IPatronRepository,
) *UserService {
	return &UserService{
		ruser: ruser,
		radmn: radmn,
		rcare: rcare,
		rmemb: rmemb,
		rorga: rorga,
		rpart: rpart,
		// rpatr: rpatr,
	}
}

func (s *UserService) Create(req requests.CreateUser) (*responses.User, error) {
	var ures *responses.User

	u := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	u, err := s.ruser.Create(*u)
	if err != nil {
		return nil, err
	}

	umar, _ := json.Marshal(u)
	err = json.Unmarshal(umar, &ures)
	if err != nil {
		return nil, err
	}

	return ures, err
}

func (s *UserService) FindByID(uid uuid.UUID) (*responses.User, error) {
	u, err := s.ruser.FindByID(uid)
	if err != nil {
		return nil, err
	}

	ures, err := u.ToResponse()
	if err != nil {
		return nil, err
	}

	return ures, err
}

func (s *UserService) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	users, err := s.ruser.FindAll(p)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) Delete(uid uuid.UUID) error {
	u := models.User{
		Model: base.Model{ID: uid},
	}

	err := s.ruser.Delete(u)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Update(
	req requests.UpdateUser,
	uid uuid.UUID,
) (*responses.User, error) {
	u, err := s.ruser.FindByID(uid)
	if err != nil {
		return nil, err
	}

	u, err = req.ToModel(*u, consttypes.UR_USER)
	if err != nil {
		return nil, err
	}

	u, err = s.ruser.Update(*u)
	if err != nil {
		return nil, err
	}

	ures, err := u.ToResponse()
	if err != nil {
		return nil, err
	}

	return ures, err
}

func (s *UserService) GetUserName(uid uuid.UUID) (string, error) {
	var name string

	user, err := s.ruser.FindByID(uid)
	if err != nil {
		return "", err
	}

	switch user.Role {
	case consttypes.UR_ADMIN:
		a, err := s.radmn.FindByUserID(uid)
		if err != nil {
			return "", err
		}

		name = fmt.Sprintf("%s %s", a.FirstName, a.LastName)
	case consttypes.UR_CAREGIVER:
		c, err := s.rcare.FindByUserID(uid)
		if err != nil {
			return "", err
		}

		name = fmt.Sprintf("%s %s", c.FirstName, c.LastName)
	case consttypes.UR_MEMBER:
		m, err := s.rmemb.FindByUserID(uid)
		if err != nil {
			return "", err
		}

		name = fmt.Sprintf("%s %s", m.FirstName, m.LastName)
	case consttypes.UR_ORGANIZATION:
		o, err := s.rorga.FindByUserID(uid)
		if err != nil {
			return "", err
		}

		name = o.Name
	case consttypes.UR_PARTNER:
		p, err := s.rpart.FindByUserID(uid)
		if err != nil {
			return "", err
		}

		name = p.Name
	case consttypes.UR_PATRON:
		// p, err := s.rpatr.FindByUserID(uid) // TODO - update this with patron data
	default:
		return "", utresponse.ErrUserInvalidRole
	}

	return name, nil
}
