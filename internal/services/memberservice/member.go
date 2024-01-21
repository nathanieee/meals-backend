package memberservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
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
		membrepo  memberrepo.IMemberRepository
		userrepo  userrepo.IUserRepository
		crgvrrepo caregiverrepo.ICaregiverRepository
		allgrepo  allergyrepo.IAllergyRepository
		illrepo   illnessrepo.IIllnessRepository
		orgrepo   organizationrepo.OrganizationRepository
	}

	IMemberService interface {
		Create(req requests.CreateMemberRequest) (*responses.MemberResponse, error)
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*responses.MemberResponse, error)
		Delete(id uuid.UUID) error
	}
)

func NewMemberService(
	membrepo memberrepo.IMemberRepository,
	userrepo userrepo.IUserRepository,
	crgvrrepo caregiverrepo.ICaregiverRepository,
	allgrepo allergyrepo.IAllergyRepository,
	illrepo illnessrepo.IIllnessRepository,
	orgrepo organizationrepo.OrganizationRepository,
) *MemberService {
	return &MemberService{
		membrepo:  membrepo,
		userrepo:  userrepo,
		crgvrrepo: crgvrrepo,
		allgrepo:  allgrepo,
		illrepo:   illrepo,
		orgrepo:   orgrepo,
	}
}

func (mes *MemberService) Create(req requests.CreateMemberRequest) (*responses.MemberResponse, error) {
	var illnesses []*models.MemberIllness
	var allergies []*models.MemberAllergy
	var caregiver *models.Caregiver
	var organization *models.Organization
	var err error

	user := req.User.ToModel(consttypes.UR_MEMBER)

	// * if caregiver request is not empty, then convert it to model.
	if req.Caregiver != nil {
		caregiver = req.Caregiver.ToModel()
	}

	// * check the organization id and assign it to the object.
	if req.OrganizationID != nil {
		organization, err = mes.orgrepo.FindByID(*req.OrganizationID)
		if err != nil {
			utlogger.LogError(err)
			return nil, err
		}
	}

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		illness, err := mes.illrepo.FindByID(ill)
		if err != nil {
			utlogger.LogError(err)
			return nil, err
		}

		millness := illness.ToMemberIllness()

		illnesses = append(illnesses, millness)
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		allergy, err := mes.allgrepo.FindByID(all)
		if err != nil {
			utlogger.LogError(err)
			return nil, err
		}

		mallergy := allergy.ToMemberAllergy()

		allergies = append(allergies, mallergy)
	}

	member := req.ToModel(*user, *caregiver, allergies, illnesses, organization)
	member, err = mes.membrepo.Create(*member)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	mres := member.ToResponse()

	return mres, nil
}

func (mes *MemberService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	members, err := mes.membrepo.FindAll(preq)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return members, nil
}

func (mes *MemberService) FindByID(id uuid.UUID) (*responses.MemberResponse, error) {
	member, err := mes.membrepo.FindByID(id)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	mres := member.ToResponse()

	return mres, nil
}

func (mes *MemberService) Delete(id uuid.UUID) error {
	member := models.Member{
		Model: helper.Model{ID: id},
	}

	err := mes.membrepo.Delete(member)
	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}
