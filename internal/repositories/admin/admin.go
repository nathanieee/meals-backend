package admnrepository

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

type (
	AdminRepository struct {
		db *gorm.DB
	}

	IAdminRepository interface {
		Create(a *models.Admin) (*models.Admin, error)
		Update(a models.Admin, aid uuid.UUID) (*models.Admin, error)
		FindAll(p models.Pagination) (*models.Pagination, error)
		FindByID(aid uuid.UUID) (*responses.AdminResponse, error)
		Delete(a models.Admin) error
	}
)

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) preload(db *gorm.DB) *gorm.DB {
	return db.
		Preload(clause.Associations)
}

func (r *AdminRepository) Create(a *models.Admin) (*models.Admin, error) {
	err := r.db.Create(a).Error
	if err != nil {
		return nil, err
	}

	return a, err
}

func (r *AdminRepository) Update(a models.Admin, aid uuid.UUID) (*models.Admin, error) {
	err := r.db.Model(&a).Where("id = ?", aid).Updates(a).Error
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *AdminRepository) FindAll(p models.Pagination) (*models.Pagination, error) {
	var a []models.Admin
	var ares []responses.AdminResponse

	result := r.db.Model(&a).Select("id, user_id, first_name, last_name, gender, date_of_birth, created_at, updated_at")

	if p.Search != "" {
		result = result.
			Where("first_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search)).
			Or("last_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search))
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.Where("date(created_at) between ? and ?", p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT), p.Filter.CreatedTo.Format(consttypes.DATEFORMAT))
	}

	result = result.Group("id").Scopes(pagination.Paginate(&a, &p, result)).Find(&ares)

	if result.Error != nil {
		return nil, result.Error
	}

	p.Data = ares
	return &p, nil
}

func (r *AdminRepository) FindByID(aid uuid.UUID) (*responses.AdminResponse, error) {
	var ares *responses.AdminResponse
	err := r.db.Model(&models.Admin{}).Select("id, user_id, first_name, last_name, gender, date_of_birth, created_at, updated_at").First(&ares, aid).Error
	if err != nil {
		return nil, err
	}

	return ares, nil
}

func (r *AdminRepository) Delete(a models.Admin) error {
	err := r.db.Delete(&a).Error
	if err != nil {
		return err
	}

	return nil
}
