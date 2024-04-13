package memberservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/base"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/organizationrepo"
	"project-skbackend/internal/repositories/userrepo"

	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	MemberService struct {
		rmemb memberrepo.IMemberRepository
		ruser userrepo.IUserRepository
		rcare caregiverrepo.ICaregiverRepository
		rall  allergyrepo.IAllergyRepository
		rill  illnessrepo.IIllnessRepository
		rorg  organizationrepo.OrganizationRepository
	}

	IMemberService interface {
		Create(req requests.CreateMember) (*responses.Member, error)
		Read() ([]*responses.Member, error)
		Update(id uuid.UUID, req requests.UpdateMember) (*responses.Member, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*responses.Member, error)
	}
)

func NewMemberService(
	rmemb memberrepo.IMemberRepository,
	ruser userrepo.IUserRepository,
	rcare caregiverrepo.ICaregiverRepository,
	rall allergyrepo.IAllergyRepository,
	rill illnessrepo.IIllnessRepository,
	rorg organizationrepo.OrganizationRepository,
) *MemberService {
	return &MemberService{
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
		return nil, err
	}

	// * if caregiver request is not empty, then convert it to model.
	if req.Caregiver != nil {
		caregiver, err = req.Caregiver.ToModel()
		if err != nil {
			return nil, err
		}
	}

	// * check the organization id and assign it to the object.
	if req.OrganizationID != nil {
		organization, err = s.rorg.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, err
		}
	}

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		illness, err := s.rill.FindByID(*ill)
		if err != nil {
			return nil, err
		}

		millness := illness.ToMemberIllness()

		illnesses = append(illnesses, millness)
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		allergy, err := s.rall.FindByID(*all)
		if err != nil {
			return nil, err
		}

		mallergy := allergy.ToMemberAllergy()

		allergies = append(allergies, mallergy)
	}

	member, err := req.ToModel(*user, caregiver, allergies, illnesses, organization)
	if err != nil {
		return nil, err
	}

	member, err = s.rmemb.Create(*member)
	if err != nil {
		return nil, err
	}

	mres := member.ToResponse()

	return mres, nil
}

func (s *MemberService) Read() ([]*responses.Member, error) {
	members, err := s.rmemb.Read()
	if err != nil {
		return nil, err
	}

	mereses := make([]*responses.Member, 0, len(members))
	for _, member := range members {
		meres := member.ToResponse()
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

	member, err := s.rmemb.FindByID(id)
	if err != nil {
		return nil, err
	}

	user, err := req.User.ToModel(member.User, consttypes.UR_MEMBER)
	if err != nil {
		return nil, err
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

			caregiver, err = s.rcare.FindByID(*member.CaregiverID)
			if err != nil {
				return nil, err
			}
		}

		caregiver, err = req.Caregiver.ToModel(caregiver)
		if err != nil {
			return nil, err
		}
	}

	// * check the organization id and assign it to the object.
	if req.OrganizationID != nil {
		organization, err = s.rorg.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, err
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
			illness, err := s.rill.FindByID(*ill)
			if err != nil {
				return nil, err
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
			allergy, err := s.rall.FindByID(*all)
			if err != nil {
				return nil, err
			}

			mallergy := allergy.ToMemberAllergy()

			allergies = append(allergies, mallergy)
		}
	}

	// * copy the request to the member model.
	member, err = req.ToModel(*member, *user, caregiver, allergies, illnesses, organization)
	if err != nil {
		return nil, err
	}

	member, err = s.rmemb.Update(*member)
	if err != nil {
		return nil, err
	}

	mres := member.ToResponse()

	return mres, nil
}

func (s *MemberService) Delete(id uuid.UUID) error {
	member := models.Member{
		Model: base.Model{ID: id},
	}

	err := s.rmemb.Delete(member)
	if err != nil {
		return err
	}

	return nil
}

func (s *MemberService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	members, err := s.rmemb.FindAll(preq)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (s *MemberService) FindByID(id uuid.UUID) (*responses.Member, error) {
	member, err := s.rmemb.FindByID(id)
	if err != nil {
		return nil, err
	}

	mres := member.ToResponse()

	return mres, nil
}
