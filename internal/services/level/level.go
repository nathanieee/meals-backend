package level

import (
	"encoding/json"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories"

	"gorm.io/gorm"
)

type LevelService struct {
	lr repositories.ILevelRepo
}

func NewLevelService(lr repositories.ILevelRepo) *LevelService {
	return &LevelService{lr: lr}
}

func (ls *LevelService) CreateLevel(req requests.CreateLevelRequest) (*responses.LevelResponse, error) {
	var lres *responses.LevelResponse

	l := &models.Level{
		Name: req.Name,
	}

	l, err := ls.lr.Store(l)
	if err != nil {
		return nil, err
	}

	marshaledLevel, _ := json.Marshal(l)
	err = json.Unmarshal(marshaledLevel, &lres)
	if err != nil {
		return nil, err
	}

	return lres, err
}

func (ls *LevelService) GetLevel(lid uint) (*responses.LevelResponse, error) {
	l, err := ls.lr.FindByID(lid)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (ls *LevelService) GetLevels(p models.Pagination) (*models.Pagination, error) {
	levels, err := ls.lr.FindAll(p)
	if err != nil {
		return nil, err
	}

	return levels, nil
}

func (ls *LevelService) DeleteLevel(lid uint) error {
	l := models.Level{
		Model: gorm.Model{ID: lid},
	}

	err := ls.lr.DeleteLevel(l)
	if err != nil {
		return err
	}

	return nil
}
