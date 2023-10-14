package member

import (
	"encoding/json"
	"math"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type MemberService struct {
	mer repositories.IMemberRepo
	ur  repositories.IUserRepo
	cgr repositories.ICaregiverRepo
	alr repositories.IAllergyRepo
}

func NewMemberService(mer repositories.IMemberRepo, ur repositories.IUserRepo) *MemberService {
	return &MemberService{mer: mer}
}

func (mes *MemberService) Create(req requests.CreateMemberRequest) (*responses.MemberResponse, error) {
	var meres *responses.MemberResponse
	var user, cgruser *models.User
	var caregiver *models.Caregiver
	// var illnesses []*models.MemberIllness // TODO - assign illness later
	var allergies []*models.MemberAllergy

	ures, err := mes.ur.FindByEmail(req.User.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user = &models.User{
				Email:    req.User.Email,
				Password: req.User.Password,
			}

			user, err = mes.ur.Create(user)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if !ures.IsEmpty() {
		err = copier.Copy(&user, &ures)
		if err != nil {
			return nil, err
		}
	}

	cgres, err := mes.cgr.FindByEmail(req.Caregiver.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cgruser = &models.User{
				Email:    req.Caregiver.Email,
				Password: req.Caregiver.Password,
			}

			cgruser, err = mes.ur.Create(cgruser)
			if err != nil {
				return nil, err
			}

			caregiver = &models.Caregiver{
				UserID:      cgruser.ID,
				User:        *cgruser,
				Gender:      req.Caregiver.Gender,
				FirstName:   req.Caregiver.FirstName,
				LastName:    req.Caregiver.LastName,
				DateOfBirth: req.Caregiver.DateOfBirth,
			}

			caregiver, err = mes.cgr.Create(caregiver)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if !cgres.IsEmpty() {
		err = copier.Copy(&caregiver, &cgres)
		if err != nil {
			return nil, err
		}
	}

	for _, alid := range req.AllergyID {
		var allergy models.Allergy
		alres, err := mes.alr.FindByID(alid)
		if err != nil {
			return nil, err
		}

		if !alres.IsEmpty() {
			err = copier.Copy(&allergy, &alres)
			if err != nil {
				return nil, err
			}
		}

		mallergy := models.MemberAllergy{
			AllergyID: allergy.ID,
			Allergy:   allergy,
		}

		allergies = append(allergies, &mallergy)
	}

	member := &models.Member{
		UserID:      user.ID,
		User:        *user,
		CaregiverID: caregiver.ID,
		Caregiver:   caregiver,
		Allergy:     allergies, // TODO - assign illness here
		Height:      req.Height,
		Weight:      req.Weight,
		BMI:         req.Weight / math.Pow(2, req.Height),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
	}

	member, err = mes.mer.Create(member)
	if err != nil {
		return nil, err
	}

	marm, _ := json.Marshal(member)
	err = json.Unmarshal(marm, &meres)
	if err != nil {
		return nil, err
	}

	return meres, nil
}
