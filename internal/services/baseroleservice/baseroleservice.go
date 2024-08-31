package baseroleservice

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utrole"
)

type (
	BaseRoleService struct {
		rmemb memberrepo.IMemberRepository
	}

	IBaseRoleService interface {
		GetMemberByBaseRole(roleres responses.BaseRole) (*models.Member, error)
	}
)

func NewBaseRoleService(
	rmemb memberrepo.IMemberRepository,
) *BaseRoleService {
	return &BaseRoleService{
		rmemb: rmemb,
	}
}

func (s *BaseRoleService) GetMemberByBaseRole(roleres responses.BaseRole) (*models.Member, error) {
	var (
		m   *models.Member
		err error
	)

	rid, rtype, ok := utrole.CartRoleCheck(roleres)
	if !ok {
		return nil, consttypes.ErrUserInvalidRole
	}

	if rtype == consttypes.UR_CAREGIVER {
		m, err = s.rmemb.GetByCaregiverID(rid)
		if err != nil {
			return nil, consttypes.ErrMemberNotFound
		}
	} else if rtype == consttypes.UR_MEMBER {
		m, err = s.rmemb.GetByID(rid)
		if err != nil {
			return nil, consttypes.ErrMemberNotFound
		}
	}

	return m, nil
}
