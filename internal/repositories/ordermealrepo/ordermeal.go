package ordermealrepo

import (
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/paginationrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	SELECTED_FIELDS = `
		id,
		order_id,
		meal_id,
		partner_id,
		quantity,
		created_at,
		updated_at
	`
)

type (
	OrderMealRepository struct {
		db  *gorm.DB
		cfg *configs.Config
	}

	IOrderMealRepository interface {
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
	}
)

func NewOrderMealRepository(
	db *gorm.DB,
	cfg *configs.Config,
) *OrderMealRepository {
	return &OrderMealRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *OrderMealRepository) omit() *gorm.DB {
	return r.db.Omit(
		"Meal",
		"Member",
	)
}

func (r *OrderMealRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Meal.Images.Image").
		Preload("Partner.User.Image.Image")
}

func (r *OrderMealRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		oms      []models.OrderMeal
		omresses []responses.OrderMeal
	)

	result := r.
		preload().
		Debug().
		Model(&oms).
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

	if p.Filter.Partner.ID != nil {
		result = result.
			Where("partner_id = ?",
				p.Filter.Partner.ID,
			)
	}

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&oms, &p, result)).
		Find(&oms)

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&omresses, &oms, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = omresses
	return &p, nil
}
