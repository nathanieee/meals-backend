package level

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/pagination"
	"project-skbackend/packages/consttypes"

	"gorm.io/gorm"
)

type LevelRepo struct {
	db *gorm.DB
}

func NewLevelRepo(db *gorm.DB) *LevelRepo {
	return &LevelRepo{db: db}
}

func (lr *LevelRepo) Store(l *models.Level) (*models.Level, error) {
	err := lr.db.Create(l).Error
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (lr *LevelRepo) Update(l models.Level, lid uint) (*models.Level, error) {
	err := lr.db.Model(&l).Where("id = ?", lid).Updates(l).Error
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (lr *LevelRepo) FindAll(p models.Pagination) (*models.Pagination, error) {
	var levels []models.Level
	var levelsResponse []responses.LevelResponse

	result := lr.db.Model(&levels)

	if p.Search != "" {
		result = result.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Search))
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.Where("date(levels.created_at) between ? and ?", p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT), p.Filter.CreatedTo.Format(consttypes.DATEFORMAT))
	}

	result = result.Group("levels.id").Scopes(pagination.Paginate(&levels, &p, result)).Find(&levelsResponse)

	if result.Error != nil {
		return &p, result.Error
	}

	p.Data = levelsResponse
	return &p, nil
}

func (lr *LevelRepo) FindByID(lid uint) (*responses.LevelResponse, error) {
	var lres *responses.LevelResponse
	err := lr.db.Model(&models.Level{}).Group("levels.id").Where("levels.deleted_at is null").First(&lres, lid).Error

	if err != nil {
		return nil, err
	}

	return lres, err
}

func (lr *LevelRepo) DeleteLevel(l models.Level) error {
	err := lr.db.Unscoped().Delete(&l).Error
	if err != nil {
		return err
	}

	return nil
}
