package adminrepo

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
		first_name,
		last_name,
		gender,
		date_of_birth,
		created_at,
		updated_at
	`
)

type (
	AdminRepository struct {
		db *gorm.DB
	}

	IAdminRepository interface {
		Create(a models.Admin) (*models.Admin, error)
		Read() ([]*models.Admin, error)
		Update(a models.Admin) (*models.Admin, error)
		Delete(a models.Admin) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*models.Admin, error)
		GetByEmail(email string) (*models.Admin, error)
		GetByUserID(uid uuid.UUID) (*models.Admin, error)
	}
)

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) omit() *gorm.DB {
	return r.db.Omit(
		"",
	)
}

func (r *AdminRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("User.Address.AddressDetail").
		Preload("User.Image.Image")
}

func (r *AdminRepository) Create(a models.Admin) (*models.Admin, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&a).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	anew, err := r.GetByID(a.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return anew, err
}

func (r *AdminRepository) Read() ([]*models.Admin, error) {
	var (
		a []*models.Admin
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&a).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return a, nil
}

func (r *AdminRepository) Update(a models.Admin) (*models.Admin, error) {
	err := r.db.
		Save(&a).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	anew, err := r.GetByID(a.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return anew, nil
}

func (r *AdminRepository) Delete(a models.Admin) error {
	err := r.db.
		Delete(&a).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *AdminRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		a    []models.Admin
		ares []responses.Admin
	)

	result := r.
		preload().
		Model(&a).
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
			Where("date(created_at) between ? and ?",
				p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT),
				p.Filter.CreatedTo.Format(consttypes.DATEFORMAT),
			)
	}

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&a, &p, result)).
		Find(&a)

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&ares, &a, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = ares
	return &p, nil
}

func (r *AdminRepository) GetByID(id uuid.UUID) (*models.Admin, error) {
	var (
		a *models.Admin
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Admin{Model: base.Model{ID: id}}).
		First(&a).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return a, nil
}

func (r *AdminRepository) GetByEmail(email string) (*models.Admin, error) {
	var (
		a *models.Admin
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(`
			admins.user_id IN (
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
		First(&a).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return a, nil
}

func (r *AdminRepository) GetByUserID(uid uuid.UUID) (*models.Admin, error) {
	var (
		a *models.Admin
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Admin{UserID: uid}, uid).
		First(&a).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return a, nil
}
