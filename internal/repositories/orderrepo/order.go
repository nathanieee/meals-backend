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

		// * this is used by cron service for automation
		UpdateAutomaticallyCancelled() error
		UpdateAutomaticallyPickedUp() error
		UpdateAutomaticallyOutForDelivery() error
		UpdateAutomaticallyDelivered() error
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

func (r *OrderRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Member.User.Image.Image").
		Preload("Member.User.Address").
		Preload("Member.Caregiver.User.Image.Image").
		Preload("Member.Caregiver.User.Address").
		Preload("Member.Organization").
		Preload("Member.Allergies.Allergy").
		Preload("Member.Illnesses.Illness").
		Preload("Meals.Meal").
		Preload("History.User.Image.Image").
		Preload("History.User.Address")
}

func (r *OrderRepository) omit() *gorm.DB {
	return r.db.Omit(
		"Member",
		"Meals.Meal",
	)
}

func (r *OrderRepository) Create(o models.Order) (*models.Order, error) {
	err := r.
		omit().
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

func (r *OrderRepository) UpdateAutomaticallyCancelled() error {
	// * get the buffer time before automatically cancelling the order
	var (
		buffer     = time.Duration(r.oacb) * time.Minute
		buffertime = consttypes.TimeNow().Add(-buffer).Format(consttypes.DATETIMEHOURMINUTESFORMAT)
		orders     []models.Order
		admin      *models.User
		status     = consttypes.OS_CANCELLED
		trigger    = []consttypes.OrderStatus{
			consttypes.OS_PLACED,
		}
	)

	admin, err := r.getAdmin()
	if err != nil {
		utlogger.Error(err)
		return err
	}

	err = r.db.
		Where("status IN ?", trigger).
		Where("TO_CHAR(created_at, '%Y-%m-%d %H:%i:00') <= ?", buffertime).
		Find(&orders).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		utlogger.Error(err)
		return err
	}

	for _, order := range orders {
		// * append the history to be cancelled
		order.History = append(order.History, models.OrderHistory{
			UserID:      admin.ID,
			User:        *admin,
			Status:      status,
			Description: consttypes.NewOrderHistoryDescription(status, admin.Email),
		})

		// * update the order main table
		err := r.db.
			Model(&order).
			Update("status", status).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			utlogger.Error(err)
			return err
		}
	}

	return nil
}

func (r *OrderRepository) UpdateAutomaticallyPickedUp() error {
	// * get the buffer time before automatically pick up the order
	var (
		buffer     = time.Duration(r.oabpu) * time.Minute
		buffertime = consttypes.TimeNow().Add(-buffer).Format(consttypes.DATETIMEHOURMINUTESFORMAT)
		orders     []models.Order
		admin      *models.User
		status     = consttypes.OS_PICKED_UP
		trigger    = []consttypes.OrderStatus{
			consttypes.OS_PREPARED,
		}
	)

	admin, err := r.getAdmin()
	if err != nil {
		utlogger.Error(err)
		return err
	}

	err = r.db.
		Where("status IN ?", trigger).
		Where("TO_CHAR(created_at, '%Y-%m-%d %H:%i:00') <= ?", buffertime).
		Find(&orders).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		utlogger.Error(err)
		return err
	}

	for _, order := range orders {
		// * append the history to be picked up
		order.History = append(order.History, models.OrderHistory{
			UserID:      admin.ID,
			User:        *admin,
			Status:      status,
			Description: consttypes.NewOrderHistoryDescription(status, admin.Email),
		})

		// * update the order main table
		err := r.db.
			Model(&order).
			Update("status", status).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			utlogger.Error(err)
			return err
		}
	}

	return nil
}

func (r *OrderRepository) UpdateAutomaticallyOutForDelivery() error {
	// * get the buffer time before automatically delivering the order
	var (
		buffer     = time.Duration(r.oabpu) * time.Minute
		buffertime = consttypes.TimeNow().Add(-buffer).Format(consttypes.DATETIMEHOURMINUTESFORMAT)
		orders     []models.Order
		admin      *models.User
		status     = consttypes.OS_OUT_FOR_DELIVERY
		trigger    = []consttypes.OrderStatus{
			consttypes.OS_PICKED_UP,
		}
	)

	admin, err := r.getAdmin()
	if err != nil {
		utlogger.Error(err)
		return err
	}

	err = r.db.
		Where("status IN ?", trigger).
		Where("TO_CHAR(created_at, '%Y-%m-%d %H:%i:00') <= ?", buffertime).
		Find(&orders).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		utlogger.Error(err)
		return err
	}

	for _, order := range orders {
		// * append the history to be delivered
		order.History = append(order.History, models.OrderHistory{
			UserID:      admin.ID,
			User:        *admin,
			Status:      status,
			Description: consttypes.NewOrderHistoryDescription(status, admin.Email),
		})

		// * update the order main table
		err := r.db.
			Model(&order).
			Update("status", status).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			utlogger.Error(err)
			return err
		}
	}

	return nil
}

func (r *OrderRepository) UpdateAutomaticallyDelivered() error {
	// * get the buffer time before the order being automatically delivered
	var (
		buffer     = time.Duration(r.oabpu) * time.Minute
		buffertime = consttypes.TimeNow().Add(-buffer).Format(consttypes.DATETIMEHOURMINUTESFORMAT)
		orders     []models.Order
		admin      *models.User
		status     = consttypes.OS_DELIVERED
		trigger    = []consttypes.OrderStatus{
			consttypes.OS_OUT_FOR_DELIVERY,
		}
	)

	admin, err := r.getAdmin()
	if err != nil {
		utlogger.Error(err)
		return err
	}

	err = r.db.
		Where("status IN ?", trigger).
		Where("TO_CHAR(created_at, '%Y-%m-%d %H:%i:00') <= ?", buffertime).
		Find(&orders).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		utlogger.Error(err)
		return err
	}

	for _, order := range orders {
		// * append the history to be delivered
		order.History = append(order.History, models.OrderHistory{
			UserID:      admin.ID,
			User:        *admin,
			Status:      status,
			Description: consttypes.NewOrderHistoryDescription(status, admin.Email),
		})

		// * update the order main table
		err := r.db.
			Model(&order).
			Update("status", status).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			utlogger.Error(err)
			return err
		}
	}

	return nil
}
