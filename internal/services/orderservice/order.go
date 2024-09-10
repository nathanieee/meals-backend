package orderservice

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/cartrepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/orderrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/baseroleservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"
	"project-skbackend/packages/utils/utslice"

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
		rcart cartrepo.ICartRepository
		rpart partnerrepo.IPartnerRepository

		sbsrl baseroleservice.IBaseRoleService

		maxord int
	}

	IOrderService interface {
		Create(req requests.CreateOrder, useroderid uuid.UUID) (*responses.Order, error)
		Read() ([]*responses.Order, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*responses.Order, error)

		GetMemberRemainingOrder(uid uuid.UUID) (*responses.OrderRemaining, error)
		FindByRoleRes(roleres responses.BaseRole) ([]*responses.Order, error)
	}
)

func NewOrderService(
	cfg *configs.Config,
	rord orderrepo.IOrderRepository,
	rmeal mealrepo.IMealRepository,
	rmemb memberrepo.IMemberRepository,
	ruser userrepo.IUserRepository,
	rcare caregiverrepo.ICaregiverRepository,
	rcart cartrepo.ICartRepository,
	rpart partnerrepo.IPartnerRepository,
	sbsrl baseroleservice.IBaseRoleService,
) *OrderService {
	return &OrderService{
		rord:  rord,
		rmeal: rmeal,
		rmemb: rmemb,
		ruser: ruser,
		rcare: rcare,
		rcart: rcart,
		rpart: rpart,

		sbsrl: sbsrl,

		maxord: cfg.OrderMax.Member,
	}
}

func (s *OrderService) Create(req requests.CreateOrder, useroderid uuid.UUID) (*responses.Order, error) {
	// * retrieves the member and user order based on the provided useroderid
	member, userorder, err := s.getMemberAndUserOrder(useroderid)
	if err != nil {
		return nil, err
	}

	// * processes the cart items and calculates the total quantity
	omeals, partner, qty, err := s.processCarts(req.CartIDs)
	if err != nil {
		return nil, err
	}

	// * checks if the daily order limit has been reached
	_, err = s.checkDailyOrderLimit(member.ID, qty)
	if err != nil {
		return nil, err
	}

	// * converts the request to an order model
	order, err := req.ToModel(*member, *userorder, omeals, *partner)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	// * creates the order in the repository
	order, err = s.rord.Create(*order)
	if err != nil {
		return nil, consttypes.ErrFailedToCreateOrder
	}

	// * converts the order model to a response
	ordres, err := order.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	// * delete the cart after processing
	err = s.rcart.DeleteByIDs(req.CartIDs)
	if err != nil {
		return nil, err
	}

	return ordres, nil
}

// * retrieves the member and user order based on the provided useroderid
func (s *OrderService) getMemberAndUserOrder(useroderid uuid.UUID) (*models.Member, *models.User, error) {
	userorder, err := s.ruser.GetByID(useroderid)
	if err != nil {
		return nil, nil, err
	}

	member, err := s.getMemberByUserID(useroderid)
	if err != nil {
		return nil, nil, err
	}

	return member, userorder, nil
}

// * processes the cart items and calculates the total quantity
func (s *OrderService) processCarts(cartIDs []uuid.UUID) ([]models.OrderMeal, *models.Partner, int, error) {
	var (
		omeals []models.OrderMeal
		qty    int
		pids   []uuid.UUID
	)

	for _, cid := range cartIDs {
		cart, err := s.rcart.GetByID(cid)
		if err != nil {
			return nil, nil, 0, consttypes.ErrCartNotFound
		}

		// * append all of the cart partner ids
		pids = append(pids, cart.PartnerID)

		omeal := models.NewCreateOrderMeals(cart.Meal, cart.Quantity)
		omeals = append(omeals, *omeal)
		qty += cart.Quantity
	}

	// * check if the order has different partner in 1 order
	if isdiffpartner := utslice.HasDifferentElements(pids); isdiffpartner {
		return nil, nil, 0, consttypes.ErrOrderShouldBeSamePartner
	}

	// * get the partner
	partner, err := s.rpart.GetByID(pids[0])
	if err != nil {
		return nil, nil, 0, consttypes.ErrPartnerNotFound
	}

	return omeals, partner, qty, nil
}

// * checks if the daily order limit has been reached
func (s *OrderService) checkDailyOrderLimit(mid uuid.UUID, qty int) (int, error) {
	dailyorder, err := s.rord.GetMemberDailyOrder(mid)
	if err != nil {
		return 0, consttypes.ErrFailedToGetDailyOrder
	}

	// * add the order to the daily order
	qty += dailyorder
	// * if the quantity of the order is greater than the
	// * daily order limit then block the user to order more
	if qty > s.maxord {
		return 0, consttypes.ErrDailyMaxOrderReached(s.maxord)
	}

	return qty, nil
}

func (s *OrderService) getMemberByUserID(uid uuid.UUID) (*models.Member, error) {
	var (
		member *models.Member
	)

	// * get the user who is ordering
	userorder, err := s.ruser.GetByID(uid)
	if err != nil {
		return nil, consttypes.ErrUserNotFound
	}

	// * get the member who is ordering
	switch userorder.Role {
	case consttypes.UR_MEMBER:
		member, err = s.rmemb.GetByUserID(uid)
		if err != nil {
			return nil, consttypes.ErrMemberNotFound
		}
	case consttypes.UR_CAREGIVER:
		caregiver, err := s.rcare.GetByUserID(uid)
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

	return member, nil
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

func (s *OrderService) GetMemberRemainingOrder(uid uuid.UUID) (*responses.OrderRemaining, error) {
	member, err := s.getMemberByUserID(uid)
	if err != nil {
		return nil, err
	}

	dailyorder, err := s.rord.GetMemberDailyOrder(member.ID)
	if err != nil {
		return nil, consttypes.ErrFailedToGetDailyOrder
	}

	remorder := s.maxord - dailyorder

	remorderrer, err := responses.NewOrderRemaining(remorder)
	if err != nil {
		return nil, err
	}

	return remorderrer, nil
}

func (s *OrderService) FindByRoleRes(roleres responses.BaseRole) ([]*responses.Order, error) {
	var (
		orderreses []*responses.Order
		err        error
	)

	switch roleres.Role {
	case consttypes.UR_MEMBER:
		orderreses, err = s.FindMemberOrder(roleres)
	case consttypes.UR_PARTNER:
		orderreses, err = s.FindPartnerOrder(roleres)
	default:
		return nil, consttypes.ErrFailedToReadOrder
	}

	if err != nil {
		return nil, err
	}

	return orderreses, nil
}

func (s *OrderService) FindMemberOrder(roleres responses.BaseRole) ([]*responses.Order, error) {
	var (
		orderreses []*responses.Order
	)

	m, err := s.sbsrl.GetMemberByBaseRole(roleres)
	if err != nil {
		return nil, err
	}

	orders, err := s.rord.FindByMemberID(m.ID)
	if err != nil {
		return nil, consttypes.ErrFailedToReadOrder
	}

	for _, order := range orders {
		ordres, err := order.ToResponse()
		if err != nil {
			return nil, consttypes.ErrConvertFailed
		}

		orderreses = append(orderreses, ordres)
	}

	return orderreses, nil
}

func (s *OrderService) FindPartnerOrder(roleres responses.BaseRole) ([]*responses.Order, error) {
	var (
		orderreses []*responses.Order
	)

	p, err := s.sbsrl.GetPartnerByBaseRole(roleres)
	if err != nil {
		return nil, err
	}

	orders, err := s.rord.FindByPartnerID(p.ID)
	if err != nil {
		return nil, consttypes.ErrFailedToReadOrder
	}

	for _, order := range orders {
		ordres, err := order.ToResponse()
		if err != nil {
			return nil, consttypes.ErrConvertFailed
		}

		orderreses = append(orderreses, ordres)
	}

	return orderreses, nil
}
