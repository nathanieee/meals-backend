package meal

import "gorm.io/gorm"

type MealRepo struct {
	db *gorm.DB
}

func NewMealRepo(db *gorm.DB) *MealRepo {
	return &MealRepo{db: db}
}
