package role

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/pagination"
	"project-skbackend/packages/consttypes"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	db = db.Preload(clause.Associations)
	return &RoleRepo{db: db}
}

func (rr *RoleRepo) Store(r *models.Role) (*models.Role, error) {
	err := rr.db.Create(r).Error
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (rr *RoleRepo) Update(r models.Role, rid uint) (*models.Role, error) {
	err := rr.db.Model(&r).Where("id = ?", rid).Updates(r).Error
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (rr *RoleRepo) FindAll(p models.Pagination) (*models.Pagination, error) {
	var roles []models.Role
	var rolesResponse []responses.RoleResponse

	result := rr.db.Model(&roles)

	if p.Search != "" {
		result = result.Where("roles.name LIKE ?", fmt.Sprintf("%%%s%%", p.Search)).Or("levels.name LIKE ?", fmt.Sprintf("%%%s%%", p.Search))
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.Where("date(role.created_at) between ? and ?", p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT), p.Filter.CreatedTo.Format(consttypes.DATEFORMAT))
	}

	result = result.Group("roles.id").Scopes(pagination.Paginate(&roles, &p, result)).Find(&rolesResponse)

	if result.Error != nil {
		return &p, result.Error
	}

	p.Data = rolesResponse
	return &p, nil
}

func (rr *RoleRepo) FindByID(lid uint) (*responses.RoleResponse, error) {
	var r *responses.RoleResponse
	err := rr.db.Model(&models.Role{}).First(&r, lid).Error
	if err != nil {
		return nil, err
	}

	return r, err
}

func (rr *RoleRepo) DeleteRole(r models.Role) error {
	err := rr.db.Unscoped().Delete(&r).Error
	if err != nil {
		return err
	}

	return nil
}
