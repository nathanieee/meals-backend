package mmbrservice

import (
	"fmt"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	allgrepository "project-skbackend/internal/repositories/allergy"
	crgvrrepository "project-skbackend/internal/repositories/caregiver"
	illrepository "project-skbackend/internal/repositories/illness"
	mmbrrepository "project-skbackend/internal/repositories/member"
	orgrepository "project-skbackend/internal/repositories/organization"
	userrepository "project-skbackend/internal/repositories/user"

	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/logger"
)

type (
	MemberService struct {
		membrepo  mmbrrepository.IMemberRepository
		userrepo  userrepository.IUserRepository
		crgvrrepo crgvrrepository.ICaregiverRepository
		allgrepo  allgrepository.IAllergyRepository
		illrepo   illrepository.IIllnessRepository
		orgrepo   orgrepository.OrganizationRepository
	}

	IMemberService interface {
		Create(req requests.CreateMemberRequest) (*responses.MemberResponse, error)
	}
)

func NewMemberService(
	membrepo mmbrrepository.IMemberRepository,
	userrepo userrepository.IUserRepository,
	crgvrrepo crgvrrepository.ICaregiverRepository,
	allgrepo allgrepository.IAllergyRepository,
	illrepo illrepository.IIllnessRepository,
	orgrepo orgrepository.OrganizationRepository,
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
			logger.LogError(err)
			return nil, err
		}
	}

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		illness, err := mes.illrepo.FindByID(ill)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}

		millness := illness.ToMemberIllness()

		illnesses = append(illnesses, millness)
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		allergy, err := mes.allgrepo.FindByID(all)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}

		mallergy := allergy.ToMemberAllergy()

		allergies = append(allergies, mallergy)
	}

	member := req.ToModel(*user, *caregiver, allergies, illnesses, organization)
	member, err = mes.membrepo.Create(*member)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	fmt.Println(member)

	mres := member.ToResponse()

	return mres, nil
}
