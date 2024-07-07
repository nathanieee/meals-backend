package orderservice

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/orderrepo"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	OrderService struct {
		rord  orderrepo.IOrderRepository
		rmeal mealrepo.IMealRepository
		rmemb memberrepo.IMemberRepository
		ruser userrepo.IUserRepository
		rcare caregiverrepo.ICaregiverRepository

		maxord uint
	}

	IOrderService interface {
		Create(req requests.CreateOrder, useroderid uuid.UUID) (*responses.Order, error)
		Read() ([]*responses.Order, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*responses.Order, error)
	}
)

func NewOrderService(
	cfg configs.Config,
	rord orderrepo.IOrderRepository,
	rmeal mealrepo.IMealRepository,
	rmemb memberrepo.IMemberRepository,
	ruser userrepo.IUserRepository,
	rcare caregiverrepo.ICaregiverRepository,
) *OrderService {
	return &OrderService{
		rord:  rord,
		rmeal: rmeal,
		rmemb: rmemb,
		ruser: ruser,
		rcare: rcare,

		maxord: cfg.OrderMax.Member,
	}
}

func (s *OrderService) Create(req requests.CreateOrder, useroderid uuid.UUID) (*responses.Order, error) {
	var (
		omeals []models.OrderMeal
		member *models.Member
		err    error
		qty    uint
	)

	// TODO: get the total order for today of the member or caregiver
	for _, omeal := range req.Meals {
		meal, err := s.rmeal.GetByID(omeal.MealID)
		if err != nil {
			return nil, consttypes.ErrMealsNotFound
		}

		omeal, err := omeal.ToModel(*meal)
		if err != nil {
			return nil, consttypes.ErrConvertFailed
		}

		omeals = append(omeals, *omeal)
		qty += omeal.Quantity
	}

	if qty >= s.maxord {
		return nil, consttypes.ErrDailyMaxOrderReached(s.maxord)
	}

	userorder, err := s.ruser.GetByID(useroderid)
	if err != nil {
		return nil, consttypes.ErrUserNotFound
	}

	switch userorder.Role {
	case consttypes.UR_MEMBER:
		member, err = s.rmemb.GetByUserID(useroderid)
		if err != nil {
			return nil, consttypes.ErrMemberNotFound
		}
	case consttypes.UR_CAREGIVER:
		caregiver, err := s.rcare.GetByUserID(useroderid)
		if err != nil {
			return nil, consttypes.ErrCaregiverNotFound
		}

		member, err = s.rmemb.GetByCaregiverID(caregiver.ID)
		if err != nil {
			return nil, consttypes.ErrMemberNotFound
		}
	default:
		return nil, consttypes.ErrUserNotFound
	}

	order, err := req.ToModel(*member, *userorder, omeals)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	order, err = s.rord.Create(*order)
	if err != nil {
		return nil, consttypes.ErrFailedToCreateOrder
	}

	ordres, err := order.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return ordres, nil
}

func (s *OrderService) Read() ([]*responses.Order, error) {
	var (
		orderreses []*responses.Order
	)

	orders, err := s.rord.Read()
	if err != nil {
		return nil, consttypes.ErrFailedToReadOrder
	}

	if err := copier.CopyWithOption(&orderreses, &orders, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, consttypes.ErrConvertFailed
	}

	return orderreses, nil
}

func (s *OrderService) Delete(id uuid.UUID) error {
	order, err := s.rord.GetByID(id)
	if err != nil {
		return consttypes.ErrOrderNotFound
	}

	if err := s.rord.Delete(*order); err != nil {
		return consttypes.ErrFailedToDeleteOrder
	}

	return nil
}

func (s *OrderService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	return s.rord.FindAll(preq)
}

func (s *OrderService) GetByID(id uuid.UUID) (*responses.Order, error) {
	order, err := s.rord.GetByID(id)
	if err != nil {
		return nil, consttypes.ErrOrderNotFound
	}

	ordres, err := order.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return ordres, nil
}
