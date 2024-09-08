package baseroleservice

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utrole"
)

type (
	BaseRoleService struct {
		rmemb memberrepo.IMemberRepository
		rpart partnerrepo.IPartnerRepository
	}

	IBaseRoleService interface {
		GetMemberByBaseRole(roleres responses.BaseRole) (*models.Member, error)
		GetPartnerByBaseRole(roleres responses.BaseRole) (*models.Partner, error)
	}
)

func NewBaseRoleService(
	rmemb memberrepo.IMemberRepository,
	rpart partnerrepo.IPartnerRepository,
) *BaseRoleService {
	return &BaseRoleService{
		rmemb: rmemb,
		rpart: rpart,
	}
}

func (s *BaseRoleService) GetMemberByBaseRole(roleres responses.BaseRole) (*models.Member, error) {
	var (
		m   *models.Member
		err error
	)

	rid, rtype, ok := utrole.RoleTranslate(roleres)
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

func (s *BaseRoleService) GetPartnerByBaseRole(roleres responses.BaseRole) (*models.Partner, error) {
	var (
		p   *models.Partner
		err error
	)

	rid, rtype, ok := utrole.RoleTranslate(roleres)
	if !ok {
		return nil, consttypes.ErrUserInvalidRole
	}

	switch rtype {
	case consttypes.UR_PARTNER:
		p, err = s.rpart.GetByID(rid)
	}

	if err != nil {
		return nil, consttypes.ErrPartnerNotFound
	}

	return p, nil
}
