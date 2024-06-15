package memberrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/base"
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
		GetByID(id uuid.UUID) (*models.Member, error)
		GetByEmail(email string) (*models.Member, error)
		GetByUserID(uid uuid.UUID) (*models.Member, error)
		GetByCaregiverID(cgid uuid.UUID) (*models.Member, error)
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
		Preload("Caregiver.User.Address").
		Preload("Organization").
		Preload("Allergies.Allergy").
		Preload("Illnesses.Illness")
}

func (r *MemberRepository) Create(m models.Member) (*models.Member, error) {
	err := r.db.
		Create(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	mnew, err := r.GetByEmail(m.User.Email)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mnew, err
}

func (r *MemberRepository) Read() ([]*models.Member, error) {
	var (
		m []*models.Member
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return m, nil
}

func (r *MemberRepository) Update(m models.Member) (*models.Member, error) {
	err := r.db.
		Save(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	mnew, err := r.GetByEmail(m.User.Email)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mnew, nil
}

func (r *MemberRepository) Delete(m models.Member) error {
	err := r.db.
		Delete(&m).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *MemberRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		m    []models.Member
		mres []responses.Member
	)

	result := r.
		preload().
		Model(&m).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(r.db.
				Where("first_name LIKE ?", p.Search).
				Or("last_name LIKE ?", p.Search),
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

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&mres, &m, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = mres
	return &p, nil
}

func (r *MemberRepository) GetByID(id uuid.UUID) (*models.Member, error) {
	var (
		m *models.Member
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Member{Model: base.Model{ID: id}}).
		First(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return m, nil
}

func (r *MemberRepository) GetByEmail(email string) (*models.Member, error) {
	var (
		m *models.Member
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(`
			members.ID IN (
				SELECT 
					id 
				FROM 
					users
				WHERE
					email = ?
					AND deleted_at IS NULL
				GROUP BY 
					id
			)
		`, email).
		First(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return m, nil
}

func (r *MemberRepository) GetByUserID(uid uuid.UUID) (*models.Member, error) {
	var (
		m *models.Member
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(`
			members.ID IN (
				SELECT 
					id 
				FROM 
					users
				WHERE
					id = ?
					AND deleted_at IS NULL
				GROUP BY 
					id
			)
		`, uid).
		First(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return m, nil
}

func (r *MemberRepository) GetByCaregiverID(cgid uuid.UUID) (*models.Member, error) {
	var (
		m *models.Member
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Member{CaregiverID: &cgid}).
		First(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return m, nil
}
