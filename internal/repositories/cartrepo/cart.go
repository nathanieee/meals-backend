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
		Create(c models.Cart) (*models.Cart, error)
		Read() ([]*models.Cart, error)
		ReadWithReference(rid uuid.UUID, rtype consttypes.UserRole) ([]*models.Cart, error)
		Update(c models.Cart) (*models.Cart, error)
		Delete(c models.Cart) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*models.Cart, error)
		FindByMemberID(mid uuid.UUID) ([]*models.Cart, error)
		FindByCaregiverID(cgid uuid.UUID) ([]*models.Cart, error)
		FindByMealID(mid uuid.UUID) ([]*models.Cart, error)
		GetCartReferenceObject(cart models.Cart) (*responses.Member, *responses.Caregiver, error)
		GetCartByMealIDAndReference(mid uuid.UUID, rid uuid.UUID, rtype consttypes.UserRole) (*models.Cart, error)
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
		utlogger.Error(err)
		return nil, err
	}

	cnew, err := r.FindByID(c.ID)

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

func (r *CartRepository) ReadWithReference(rid uuid.UUID, rtype consttypes.UserRole) ([]*models.Cart, error) {
	var (
		c []*models.Cart
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Cart{ReferenceID: rid, ReferenceType: rtype}).
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

	cnew, err := r.FindByID(c.ID)

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

	for _, cart := range c {
		for _, cartres := range cres {
			if cart.ID == cartres.ID {
				mres, cgres, err := r.GetCartReferenceObject(cart)
				if err != nil {
					utlogger.Error(err)
					return nil, err
				}

				cartres.Member = mres
				cartres.Caregiver = cgres
			}
		}
	}

	p.Data = cres
	return &p, nil
}

func (r *CartRepository) FindByID(id uuid.UUID) (*models.Cart, error) {
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
		Where(&models.Cart{ReferenceID: mid, ReferenceType: consttypes.UR_MEMBER}).
		Find(&c).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return c, nil
}

func (r *CartRepository) FindByCaregiverID(cid uuid.UUID) ([]*models.Cart, error) {
	var (
		c []*models.Cart
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Cart{ReferenceID: cid, ReferenceType: consttypes.UR_CAREGIVER}).
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

func (r *CartRepository) GetCartReferenceObject(cart models.Cart) (*responses.Member, *responses.Caregiver, error) {
	var (
		cg    models.Caregiver
		m     models.Member
		cgres *responses.Caregiver
		mres  *responses.Member
	)

	switch cart.ReferenceType {
	case consttypes.UR_CAREGIVER:
		err := r.db.
			Preload(clause.Associations).
			Preload("User.Image.Image").
			Preload("User.Address").
			First(&cg, cart.ReferenceID).Error

		if err != nil {
			utlogger.Error(err)
			return nil, nil, err
		}

		cgres = cg.ToResponse()
	case consttypes.UR_MEMBER:
		err := r.db.
			Preload(clause.Associations).
			Preload("User.Image.Image").
			Preload("User.Address").
			Preload("Caregiver.User.Image.Image").
			Preload("Caregiver.User.Address").
			Preload("Organization").
			Preload("Allergies.Allergy").
			Preload("Illnesses.Illness").
			First(&m, cart.ReferenceID).Error

		if err != nil {
			utlogger.Error(err)
			return nil, nil, err
		}

		mres = m.ToResponse()
	default:
		return nil, nil, consttypes.ErrInvalidReference
	}

	return mres, cgres, nil
}

func (r *CartRepository) GetCartByMealIDAndReference(mid uuid.UUID, rid uuid.UUID, rtype consttypes.UserRole) (*models.Cart, error) {
	var (
		c *models.Cart
	)

	err := r.db.
		Select(SELECTED_FIELDS).
		Where(&models.Cart{MealID: mid, ReferenceID: rid, ReferenceType: rtype}).
		First(&c).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return c, nil
}
