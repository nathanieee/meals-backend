package caregiver

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/pagination"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CaregiverRepo struct {
	db *gorm.DB
}

func NewCaregiverRepo(db *gorm.DB) *CaregiverRepo {
	return &CaregiverRepo{db: db}
}

func (cgr *CaregiverRepo) Create(cg *models.Caregiver) (*models.Caregiver, error) {
	err := cgr.db.Create(cg).Error
	if err != nil {
		return nil, err
	}

	return cg, err
}

func (cgr *CaregiverRepo) Update(cg models.Caregiver, cgid uuid.UUID) (*models.Caregiver, error) {
	err := cgr.db.Model(&cg).Where("id = ?", cgid).Updates(cg).Error
	if err != nil {
		return nil, err
	}

	return &cg, nil
}

func (cgr *CaregiverRepo) FindAll(p models.Pagination) (*models.Pagination, error) {
	var cg []models.Caregiver
	var cgres []responses.CaregiverResponse

	result := cgr.db.Model(&cg).Select("id, user_id, gender, first_name, last_name, date_of_birth, created_at, updated_at")

	if p.Search != "" {
		result = result.
			Where("first_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search)).
			Or("last_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search))
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.Where("date(created_at) between ? and ?", p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT), p.Filter.CreatedTo.Format(consttypes.DATEFORMAT))
	}

	result = result.Group("id").Scopes(pagination.Paginate(&cg, &p, result)).Find(&cgres)

	if result.Error != nil {
		return nil, result.Error
	}

	p.Data = cgres
	return &p, nil
}

func (cgr *CaregiverRepo) FindByID(cgid uuid.UUID) (*responses.CaregiverResponse, error) {
	var cgres *responses.CaregiverResponse
	err := cgr.db.Model(&models.Caregiver{}).Select("id, user_id, gender, first_name, last_name, date_of_birth, created_at, updated_at").First(&cgres, cgid).Error
	if err != nil {
		return nil, err
	}

	return cgres, nil
}

func (cgr *CaregiverRepo) Delete(cg models.Caregiver) error {
	err := cgr.db.Delete(&cg).Error
	if err != nil {
		return err
	}

	return nil
}
