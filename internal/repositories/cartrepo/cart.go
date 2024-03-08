package cartrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/internal/repositories/paginationrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"
	"project-skbackend/packages/utils/utresponse"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	SELECTED_FIELDS = `
		id, 
		meal_id,
		reference_id,
		reference_type,
		quantity,
		created_at,
		updated_at
	`
)

type (
	CartRepository struct {
		db *gorm.DB
	}

	ICartRepository interface {
		Create(m models.Cart) (*models.Cart, error)
		Read() ([]*models.Cart, error)
		Update(m models.Cart) (*models.Cart, error)
		Delete(m models.Cart) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*models.Cart, error)
		FindByMemberID(mid uuid.UUID) ([]*models.Cart, error)
		FindByCaregiverID(cgid uuid.UUID) ([]*models.Cart, error)
		FindByMealID(mid uuid.UUID) ([]*models.Cart, error)
	}
)

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Meal.Images.Image").
		Preload("Meal.Illnesses.Illness").
		Preload("Meal.Allergies.Allergy").
		Preload("Meal.Partner.User.Address").
		Preload("Meal.Partner.User.Image.Image")
}

func (r *CartRepository) Create(c models.Cart) (*models.Cart, error) {
	err := r.db.
		Create(&c).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &c, nil
}

func (r *CartRepository) Read() ([]*models.Cart, error) {
	var c []*models.Cart

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&c).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return c, nil
}

func (r *CartRepository) Update(c models.Cart) (*models.Cart, error) {
	err := r.db.
		Save(&c).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &c, nil
}

func (r *CartRepository) Delete(c models.Cart) error {
	err := r.db.
		Delete(&c).Error

	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (r *CartRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var c []models.Cart
	var cres []responses.Cart
	var member *models.Member
	var caregiver *models.Caregiver

	result := r.
		preload().
		Model(&c).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(r.db.
				Where(&models.Cart{Meal: models.Meal{Name: p.Search}}),
			)
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
		Scopes(paginationrepo.Paginate(&c, &p, result)).
		Find(&c)

	if result.Error != nil {
		utlogger.LogError(result.Error)
		return nil, result.Error
	}

	// * copy the data from model to response
	copier.CopyWithOption(&cres, &c, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// TODO - fix this to be able to get member and caregiver
	for _, carts := range c {
		for _, response := range cres {
			if carts.ID == response.ID {
				switch carts.ReferenceType {
				case consttypes.UR_MEMBER:
					r.db.Model(&carts).Association("Member").Find(&member)
				case consttypes.UR_CAREGIVER:
					r.db.Model(&carts).Association("Caregiver").Find(&caregiver)
				default:
					return nil, utresponse.ErrInvalidReference
				}

				response.Member = member.ToResponse()
				response.Caregiver = caregiver.ToResponse()
			}
		}
	}

	p.Data = cres
	return &p, nil
}

func (r *CartRepository) FindByID(id uuid.UUID) (*models.Cart, error) {
	var c *models.Cart

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Cart{Model: helper.Model{ID: id}}).
		First(&c).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return c, nil
}

func (r *CartRepository) FindByMemberID(mid uuid.UUID) ([]*models.Cart, error) {
	var c []*models.Cart

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Cart{ReferenceID: mid, ReferenceType: consttypes.UR_MEMBER}).
		Find(&c).Error

	return nil, err
}
