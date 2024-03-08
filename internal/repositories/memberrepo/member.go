package memberrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/internal/repositories/paginationrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
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
		Read() ([]*models.Member, error)
		Update(m models.Member) (*models.Member, error)
		Delete(m models.Member) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*models.Member, error)
		FindByEmail(email string) (*models.Member, error)
	}
)

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("User.Image.Image").
		Preload("User.Address").
		Preload("Caregiver.User.Image.Image").
		Preload("Organization").
		Preload("Allergies.Allergy").
		Preload("Illnesses.Illness")
}

func (r *MemberRepository) Create(m models.Member) (*models.Member, error) {
	err := r.db.
		Create(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &m, err
}

func (r *MemberRepository) Read() ([]*models.Member, error) {
	var m []*models.Member

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return m, nil
}

func (r *MemberRepository) Update(m models.Member) (*models.Member, error) {
	err := r.db.
		Save(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &m, nil
}

func (r *MemberRepository) Delete(m models.Member) error {
	err := r.db.
		Delete(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (r *MemberRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var m []models.Member
	var mres []responses.Member

	result := r.
		preload().
		Model(&m).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(r.db.
				Where(&models.Member{FirstName: p.Search}).
				Or(&models.Member{LastName: p.Search}),
			)
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.
			Where("date(created_at) BETWEEN ? and ?",
				p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT),
				p.Filter.CreatedTo.Format(consttypes.DATEFORMAT),
			)
	}

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&m, &p, result)).
		Find(&m)

	if result.Error != nil {
		utlogger.LogError(result.Error)
		return nil, result.Error
	}

	// * copy the data from model to response
	copier.CopyWithOption(&mres, &m, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = mres
	return &p, nil
}

func (r *MemberRepository) FindByID(id uuid.UUID) (*models.Member, error) {
	var m *models.Member

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Member{Model: helper.Model{ID: id}}).
		First(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return m, nil
}

func (r *MemberRepository) FindByEmail(email string) (*models.Member, error) {
	var m *models.Member

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Member{User: models.User{Email: email}}).
		First(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return m, nil
}
