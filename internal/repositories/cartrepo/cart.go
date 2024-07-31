package cartrepo

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/base"
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
		meal_id,
		member_id,
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
		Create(c models.Cart) (*models.Cart, error)
		Read() ([]*models.Cart, error)
		Update(c models.Cart) (*models.Cart, error)
		Delete(c models.Cart) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*models.Cart, error)
		FindByMemberID(mid uuid.UUID) ([]*models.Cart, error)
		FindByMealID(mid uuid.UUID) ([]*models.Cart, error)
		GetByMealIDAndMemberID(membid uuid.UUID, mealid uuid.UUID) (*models.Cart, error)
	}
)

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) omit() *gorm.DB {
	return r.db.Omit(
		"Meal",
		"Member",
	)
}

func (r *CartRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Meal.Images.Image").
		Preload("Meal.Illnesses.Illness").
		Preload("Meal.Allergies.Allergy").
		Preload("Meal.Partner.User.Addresses.AddressDetail").
		Preload("Meal.Partner.User.Image.Image").
		Preload("Member.User.Image.Image").
		Preload("Member.User.Addresses.AddressDetail").
		Preload("Member.Caregiver.User.Image.Image").
		Preload("Member.Caregiver.User.Addresses.AddressDetail").
		Preload("Member.Organization").
		Preload("Member.Allergies.Allergy").
		Preload("Member.Illnesses.Illness")
}

func (r *CartRepository) Create(c models.Cart) (*models.Cart, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&c).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cnew, err := r.GetByID(c.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cnew, nil
}

func (r *CartRepository) Read() ([]*models.Cart, error) {
	var (
		c []*models.Cart
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&c).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return c, nil
}

func (r *CartRepository) Update(c models.Cart) (*models.Cart, error) {
	err := r.db.
		Model(&c).
		Updates(&c).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cnew, err := r.GetByID(c.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cnew, nil
}

func (r *CartRepository) Delete(c models.Cart) error {
	err := r.db.
		Delete(&c).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *CartRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		c    []models.Cart
		cres []responses.Cart
	)

	result := r.
		preload().
		Model(&c).
		Select(SELECTED_FIELDS)

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
		utlogger.Error(result.Error)
		return nil, result.Error
	}

	// * copy the data from model to response
	copier.CopyWithOption(&cres, &c, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = cres
	return &p, nil
}

func (r *CartRepository) GetByID(id uuid.UUID) (*models.Cart, error) {
	var (
		c *models.Cart
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Cart{Model: base.Model{ID: id}}).
		First(&c).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return c, nil
}

func (r *CartRepository) FindByMemberID(mid uuid.UUID) ([]*models.Cart, error) {
	var (
		c []*models.Cart
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Cart{MemberID: mid}).
		Find(&c).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return c, nil
}

func (r *CartRepository) FindByMealID(mid uuid.UUID) ([]*models.Cart, error) {
	var (
		c []*models.Cart
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Cart{MealID: mid}).
		Find(&c).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return c, nil
}

func (r *CartRepository) GetByMealIDAndMemberID(membid uuid.UUID, mealid uuid.UUID) (*models.Cart, error) {
	var (
		c *models.Cart
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Cart{MemberID: membid, MealID: mealid}).
		First(&c).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return c, nil
}
