package orderrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
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
		created_at,
		updated_at,
		member_id,
		meal_id,
		status
	`
)

type (
	OrderRepository struct {
		db *gorm.DB
	}

	IOrderRepository interface {
		Create(o models.Order) (*models.Order, error)
		Read() ([]*models.Order, error)
		Update(o models.Order) (*models.Order, error)
		Delete(o models.Order) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*models.Order, error)
		GetByMemberID(id uuid.UUID) ([]*models.Order, error)
		GetByMealID(id uuid.UUID) ([]*models.Order, error)
	}
)

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Member.User.Image.Image").
		Preload("Member.User.Address").
		Preload("Member.Caregiver.User.Image.Image").
		Preload("Member.Caregiver.User.Address").
		Preload("Member.Organization").
		Preload("Member.Allergies.Allergy").
		Preload("Member.Illnesses.Illnes").
		Preload("Meal.Images.Image").
		Preload("Meal.Illnesses.Illness").
		Preload("Meal.Allergies.Allergy").
		Preload("Meal.Partner.User.Address").
		Preload("Meal.Partner.User.Image.Imag")
}

func (r *OrderRepository) Create(o models.Order) (*models.Order, error) {
	err := r.db.
		Create(&o).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	onew, err := r.GetByID(o.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return onew, err
}

func (r *OrderRepository) Read() ([]*models.Order, error) {
	var (
		o []*models.Order
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&o).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return o, nil
}

func (r *OrderRepository) Update(o models.Order) (*models.Order, error) {
	err := r.db.
		Save(&o).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	onew, err := r.GetByID(o.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return onew, nil
}

func (r *OrderRepository) Delete(o models.Order) error {
	err := r.db.
		Delete(&o).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *OrderRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		o    []models.Order
		ores []responses.Order
	)

	result := r.
		preload().
		Model(&o).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		// TODO: add a like query
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
		Scopes(paginationrepo.Paginate(&o, &p, result)).
		Find(&ores)

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&ores, &o, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = ores
	return &p, nil
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*models.Order, error) {
	var (
		o *models.Order
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("id = ?", id).
		First(&o).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return o, nil
}

func (r *OrderRepository) GetByMealID(id uuid.UUID) ([]*models.Order, error) {
	var (
		o []*models.Order
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("meal_id = ?", id).
		Find(&o).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return o, nil
}

func (r *OrderRepository) GetByMemberID(id uuid.UUID) ([]*models.Order, error) {
	var (
		o []*models.Order
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("member_id = ?", id).
		Find(&o).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return o, nil
}
