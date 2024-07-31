package orderrepo

import (
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/paginationrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"
	"time"

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
		status
	`
)

type (
	OrderRepository struct {
		db  *gorm.DB
		cfg configs.Config

		oacb  int
		oabpu int
		oaofd int
		oad   int
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
		GetMemberDailyOrder(id uuid.UUID) (uint, error)

		// * this is used by cron service for automation
		UpdateAutomaticallyStatus(status consttypes.OrderStatus, bufferminutes int, trigger []consttypes.OrderStatus) error
	}
)

func NewOrderRepository(
	db *gorm.DB,
	cfg configs.Config,
) *OrderRepository {
	return &OrderRepository{
		db:  db,
		cfg: cfg,

		oacb:  cfg.OrderBuffer.AutomaticallyCancelled,
		oabpu: cfg.OrderBuffer.AutomaticallyBeingPickedUp,
		oaofd: cfg.OrderBuffer.AutomaticallyOutForDelivery,
		oad:   cfg.OrderBuffer.AutomaticallyDelivered,
	}
}

func (r *OrderRepository) omit() *gorm.DB {
	return r.db.Omit(
		"Member",
		"Meals.Meal",
		"Meals.Partner",
		"History.User",
	)
}

func (r *OrderRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Member.User.Image.Image").
		Preload("Member.User.Addresses.AddressDetail").
		Preload("Member.Caregiver.User.Image.Image").
		Preload("Member.Caregiver.User.Addresses.AddressDetail").
		Preload("Member.Organization").
		Preload("Member.Allergies.Allergy").
		Preload("Member.Illnesses.Illness").
		Preload("Meals.Meal").
		Preload("History.User.Image.Image").
		Preload("History.User.Addresses.AddressDetail")
}

func (r *OrderRepository) Create(o models.Order) (*models.Order, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
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

	return onew, nil
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
	err := r.
		omit().
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

func (r *OrderRepository) getAdmin() (*models.User, error) {
	var (
		admin models.User
	)

	err := r.db.
		Where("role = ?", consttypes.UR_ADMIN).
		First(&admin).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &admin, nil
}

func (r *OrderRepository) UpdateAutomaticallyStatus(status consttypes.OrderStatus, bufferminutes int, trigger []consttypes.OrderStatus) error {
	// * get the buffer time
	buffer := time.Duration(bufferminutes) * time.Minute
	buffertime := consttypes.TimeNow().Add(-buffer).Format(consttypes.DATETIMEHOURMINUTESFORMAT)

	// * find orders that meet the condition
	var orders []models.Order
	err := r.db.
		Where("status IN ?", trigger).
		Where("TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI') = ?", buffertime).
		Find(&orders).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		utlogger.Error(err)
		return err
	}

	// * get the admin user
	admin, err := r.getAdmin()
	if err != nil {
		utlogger.Error(err)
		return err
	}

	// * loop through orders and update them
	for _, order := range orders {
		// Append history
		order.History = append(order.History, models.OrderHistory{
			UserID:      admin.ID,
			User:        *admin,
			Status:      status,
			Description: consttypes.NewOrderHistoryDescription(status, admin.Email),
		})

		// * update the order status
		err := r.db.Model(&order).Updates(models.Order{Status: status}).Error
		if err != nil {
			utlogger.Error(err)
			return err
		}
	}

	return nil
}

func (r *OrderRepository) GetMemberDailyOrder(id uuid.UUID) (uint, error) {
	var (
		orders []models.Order
		qty    uint = 0
	)

	err := r.
		preload().
		Where("member_id = ?", id).
		Where("DATE(created_at) = ?::DATE", consttypes.TimeNow().Format(consttypes.DATEFORMAT)).
		Find(&orders).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		utlogger.Error(err)
		return 0, err
	}

	for _, order := range orders {
		for _, meal := range order.Meals {
			qty += meal.Quantity
		}
	}

	return qty, nil
}
