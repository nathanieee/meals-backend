package mmbrrepository

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/pagination"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils"
	"project-skbackend/packages/utils/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	SELECTED_FIELDS = `
		id,
		user_id,
		caregiver_id,
		organization_id,
		height,
		weight,
		bmi,
		first_name,
		last_name,
		gender,
		date_of_birth,
		created_at,
		updated_at
	`
)

type (
	MemberRepository struct {
		db *gorm.DB
	}

	IMemberRepository interface {
		Create(m models.Member) (*models.Member, error)
		Update(m models.Member, mid uuid.UUID) (*models.Member, error)
		FindAll(p utils.Pagination) (*utils.Pagination, error)
		FindByID(mid uuid.UUID) (*models.Member, error)
		Delete(m models.Member) error
	}
)

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) preload(db *gorm.DB) *gorm.DB {
	return db.
		Preload(clause.Associations).
		Preload("Users.UserImages.Images").
		Preload("Users.Addresses").
		Preload("Caregiver.Users.UserImages.Images").
		Preload("Organizations").
		Preload("MemberAllergies.Allergies").
		Preload("MemberIllnesses.Illnesses")
}

func (r *MemberRepository) Create(m models.Member) (*models.Member, error) {
	err := r.db.
		Create(&m).Error

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return &m, err
}

func (r *MemberRepository) Update(m models.Member, mid uuid.UUID) (*models.Member, error) {
	err := r.db.
		Model(&m).
		Where("id = ?", mid).
		Updates(m).Error

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return &m, nil
}

func (r *MemberRepository) FindAll(p utils.Pagination) (*utils.Pagination, error) {
	var m []models.Member
	var mres []responses.MemberResponse

	result := r.
		preload(r.db).
		Model(&m).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		result = result.
			Where("first_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search)).
			Or("last_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search))
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.
			Where("date(created_at) between ? and ?",
				p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT),
				p.Filter.CreatedTo.Format(consttypes.DATEFORMAT),
			)
	}

	result = result.Group("id").Scopes(pagination.Paginate(&m, &p, result)).Find(&mres)

	if result.Error != nil {
		logger.LogError(result.Error)
		return nil, result.Error
	}

	p.Data = mres
	return &p, nil
}

func (r *MemberRepository) FindByID(mid uuid.UUID) (*models.Member, error) {
	var m *models.Member

	err := r.db.
		Model(&models.Member{}).
		Select(SELECTED_FIELDS).
		First(&m, mid).Error

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return m, nil
}

func (r *MemberRepository) Delete(m models.Member) error {
	err := r.db.
		Delete(&m).Error

	if err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}
