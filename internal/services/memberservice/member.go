package memberservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/organizationrepo"
	"project-skbackend/internal/repositories/userrepo"

	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"
	"project-skbackend/packages/utils/utstring"

	"github.com/google/uuid"
)

type (
	MemberService struct {
		// * repository
		rmemb memberrepo.IMemberRepository
		ruser userrepo.IUserRepository
		rcare caregiverrepo.ICaregiverRepository
		rall  allergyrepo.IAllergyRepository
		rill  illnessrepo.IIllnessRepository
		rorg  organizationrepo.IOrganizationRepository
	}

	IMemberService interface {
		Create(req requests.CreateMember) (*responses.Member, error)
		Read() ([]*responses.Member, error)
		Update(id uuid.UUID, req requests.UpdateMember) (*responses.Member, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*responses.Member, error)
		GetByCaregiverID(cgid uuid.UUID) (*responses.Member, error)

		UpdateOwnCaregiver(mid uuid.UUID, req requests.UpdateCaregiver) (*responses.Caregiver, error)
		UpdateOwnCaregiverPassword(mid uuid.UUID, req requests.UpdatePassword) error
	}
)

func NewMemberService(
	// * repository
	rmemb memberrepo.IMemberRepository,
	ruser userrepo.IUserRepository,
	rcare caregiverrepo.ICaregiverRepository,
	rall allergyrepo.IAllergyRepository,
	rill illnessrepo.IIllnessRepository,
	rorg organizationrepo.IOrganizationRepository,
) *MemberService {
	return &MemberService{
		// * repository
		rmemb: rmemb,
		ruser: ruser,
		rcare: rcare,
		rall:  rall,
		rill:  rill,
		rorg:  rorg,
	}
}

func (s *MemberService) Create(req requests.CreateMember) (*responses.Member, error) {
	var (
		illnesses    []*models.MemberIllness
		allergies    []*models.MemberAllergy
		caregiver    *models.Caregiver
		organization *models.Organization
		err          error
	)

	user, err := req.User.ToModel(consttypes.UR_MEMBER)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	// * if caregiver request is not empty, then convert it to model.
	if req.Caregiver != nil {
		caregiver, err = req.Caregiver.FromMemberAddition()
		if err != nil {
			return nil, consttypes.ErrConvertFailed
		}
	}

	// * check the organization id and assign it to the object.
	if req.OrganizationID != nil {
		organization, err = s.rorg.GetByID(*req.OrganizationID)
		if err != nil {
			return nil, consttypes.ErrOrganizationNotFound
		}
	}

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		illness, err := s.rill.GetByID(*ill)
		if err != nil {
			return nil, consttypes.ErrIllnessNotFound
		}

		millness := illness.ToMemberIllness()

		illnesses = append(illnesses, millness)
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		allergy, err := s.rall.GetByID(*all)
		if err != nil {
			return nil, consttypes.ErrAllergiesNotFound
		}

		mallergy := allergy.ToMemberAllergy()

		allergies = append(allergies, mallergy)
	}

	member, err := req.ToModel(*user, caregiver, allergies, illnesses, organization)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	member, err = s.rmemb.Create(*member)
	if err != nil {
		return nil, consttypes.ErrFailedToCreateMember
	}

	mres, err := member.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return mres, nil
}

func (s *MemberService) Read() ([]*responses.Member, error) {
	var (
		mereses []*responses.Member
	)

	members, err := s.rmemb.Read()
	if err != nil {
		return nil, consttypes.ErrFailedToReadMembers
	}

	for _, member := range members {
		meres, err := member.ToResponse()
		if err != nil {
			return nil, consttypes.ErrConvertFailed
		}

		mereses = append(mereses, meres)
	}

	return mereses, nil
}

func (s *MemberService) Update(id uuid.UUID, req requests.UpdateMember) (*responses.Member, error) {
	var (
		illnesses    []*models.MemberIllness
		allergies    []*models.MemberAllergy
		caregiver    *models.Caregiver
		organization *models.Organization
		err          error
	)

	member, err := s.rmemb.GetByID(id)
	if err != nil {
		return nil, consttypes.ErrMemberNotFound
	}

	user, err := req.User.ToModel(member.User, consttypes.UR_MEMBER)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	// * if caregiver request is not empty, check whether the member already has one.
	// * if not, then convert it to model.
	if req.Caregiver != nil {
		if member.Caregiver != nil {
			if member.Caregiver.User.Email != req.Caregiver.User.Email {
				err := consttypes.ErrCannotChangeEmail

				utlogger.Error(err)
				return nil, err
			}

			caregiver, err = s.rcare.GetByID(*member.CaregiverID)
			if err != nil {
				return nil, consttypes.ErrCaregiverNotFound
			}
		}

		caregiver, err = req.Caregiver.ToModel(caregiver)
		if err != nil {
			return nil, consttypes.ErrConvertFailed
		}
	}

	// * check the organization id and assign it to the object.
	if req.OrganizationID != nil {
		organization, err = s.rorg.GetByID(*req.OrganizationID)
		if err != nil {
			return nil, consttypes.ErrOrganizationNotFound
		}
	}

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		var (
			found = false
		)

		for _, mill := range member.Illnesses {
			if *ill == mill.Illness.ID {
				found = true
				continue
			}
		}

		if found {
			continue
		} else {
			illness, err := s.rill.GetByID(*ill)
			if err != nil {
				return nil, consttypes.ErrIllnessNotFound
			}

			millness := illness.ToMemberIllness()

			illnesses = append(illnesses, millness)
		}
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		var (
			found = false
		)

		for _, mall := range member.Allergies {
			if *all == mall.Allergy.ID {
				found = true
				continue
			}
		}

		if found {
			continue
		} else {
			allergy, err := s.rall.GetByID(*all)
			if err != nil {
				return nil, consttypes.ErrAllergiesNotFound
			}

			mallergy := allergy.ToMemberAllergy()

			allergies = append(allergies, mallergy)
		}
	}

	// * copy the request to the member model.
	member, err = req.ToModel(*member, *user, caregiver, allergies, illnesses, organization)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	member, err = s.rmemb.Update(*member)
	if err != nil {
		return nil, consttypes.ErrFailedToUpdateMember
	}

	mres, err := member.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return mres, nil
}

func (s *MemberService) Delete(id uuid.UUID) error {
	member, err := s.rmemb.GetByID(id)
	if err != nil {
		return consttypes.ErrMemberNotFound
	}

	if err := s.rmemb.Delete(*member); err != nil {
		return consttypes.ErrFailedToDeleteMember
	}

	return nil
}

func (s *MemberService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	members, err := s.rmemb.FindAll(preq)
	if err != nil {
		return nil, consttypes.ErrFailedToFindAllMembers
	}

	return members, nil
}

func (s *MemberService) GetByID(id uuid.UUID) (*responses.Member, error) {
	member, err := s.rmemb.GetByID(id)
	if err != nil {
		return nil, consttypes.ErrMemberNotFound
	}

	mres, err := member.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return mres, nil
}

func (s *MemberService) GetByCaregiverID(cgid uuid.UUID) (*responses.Member, error) {
	member, err := s.rmemb.GetByCaregiverID(cgid)
	if err != nil {
		return nil, consttypes.ErrMemberNotFound
	}

	mres, err := member.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return mres, nil
}

func (s *MemberService) UpdateOwnCaregiver(mid uuid.UUID, req requests.UpdateCaregiver) (*responses.Caregiver, error) {
	var (
		caregiver *models.Caregiver
	)

	member, err := s.rmemb.GetByID(mid)
	if err != nil {
		return nil, consttypes.ErrMemberNotFound
	}

	if member.Caregiver != nil {
		if req.User != nil {
			if member.Caregiver.User.Email != req.User.Email && req.User.Email != "" {
				return nil, consttypes.ErrCannotChangeEmail
			}
		}

		caregiver, err = s.rcare.GetByID(*member.CaregiverID)
		if err != nil {
			return nil, consttypes.ErrCaregiverNotFound
		}
	}

	caregiver, err = req.ToModel(caregiver)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	caregiver, err = s.rcare.Update(*caregiver)
	if err != nil {
		return nil, err
	}

	cres, err := caregiver.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return cres, nil
}

func (s *MemberService) UpdateOwnCaregiverPassword(mid uuid.UUID, req requests.UpdatePassword) error {
	// * get the member by its id
	member, err := s.rmemb.GetByID(mid)
	if err != nil {
		return err
	}

	// * check if the caregiver exist or not
	if member.Caregiver == nil {
		return consttypes.ErrCaregiverNotFound
	}

	// * check if the old password is correct
	ok := utstring.CheckPasswordHash(req.OldPassword, member.Caregiver.User.Password)
	if !ok {
		return consttypes.ErrInvalidEmailOrPassword
	}

	// * update the caregiver's password
	_, err = s.ruser.UpdatePassword(member.Caregiver.User.ID, req.NewPassword)
	if err != nil {
		return err
	}

	return nil
}
