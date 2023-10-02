package member

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/pagination"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MemberRepo struct {
	db *gorm.DB
}

func NewMemberRepo(db *gorm.DB) *MemberRepo {
	db.
		Preload(clause.Associations).
		Preload("Users.UserImages.Images").
		Preload("Users.Addresses").
		Preload("Caregiver.Users.UserImages.Images").
		Preload("Organizations").
		Preload("MemberAllergies.Allergies").
		Preload("MemberIllnesses.Illnesses")

	return &MemberRepo{db: db}
}

func (mr *MemberRepo) Create(m *models.Member) (*models.Member, error) {
	err := mr.db.Create(m).Error
	if err != nil {
		return nil, err
	}

	return m, err
}

func (mr *MemberRepo) Update(m models.Member, mid uuid.UUID) (*models.Member, error) {
	err := mr.db.Model(&m).Where("id = ?", mid).Updates(m).Error
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (mr *MemberRepo) FindAll(p models.Pagination) (*models.Pagination, error) {
	var m []models.Member
	var mres []responses.MemberResponse

	result := mr.db.Model(&m).Select("id, user_id, caregiver_id, organization_id, height, weight, bmi, first_name, last_name, gender, date_of_birth, created_at, updated_at")

	if p.Search != "" {
		result = result.
			Where("first_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search)).
			Or("last_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search))
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.Where("date(created_at) between ? and ?", p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT), p.Filter.CreatedTo.Format(consttypes.DATEFORMAT))
	}

	result = result.Group("id").Scopes(pagination.Paginate(&m, &p, result)).Find(&mres)

	if result.Error != nil {
		return nil, result.Error
	}

	p.Data = mres
	return &p, nil
}

func (mr *MemberRepo) FindByID(mid uuid.UUID) (*responses.MemberResponse, error) {
	var mres *responses.MemberResponse
	err := mr.db.Model(&models.Member{}).Select("id, user_id, caregiver_id, organization_id, height, weight, bmi, first_name, last_name, gender, date_of_birth, created_at, updated_at").First(&mres, mid).Error
	if err != nil {
		return nil, err
	}

	return mres, nil
}

func (mr *MemberRepo) Delete(m models.Member) error {
	err := mr.db.Delete(&m).Error
	if err != nil {
		return err
	}

	return nil
}
