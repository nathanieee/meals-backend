package mealrepository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	MealRepository struct {
		db *gorm.DB
	}

	IMealRepository interface{}
)

func NewMealRepository(db *gorm.DB) *MealRepository {
	return &MealRepository{db: db}
}

func (r *MealRepository) preload(db *gorm.DB) *gorm.DB {
	return db.
		Preload(clause.Associations)
}
