package cartservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/cartrepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/services/baseroleservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	CartService struct {
		rcart cartrepo.ICartRepository
		rcare caregiverrepo.ICaregiverRepository
		rmemb memberrepo.IMemberRepository
		rmeal mealrepo.IMealRepository

		sbsrl baseroleservice.IBaseRoleService
	}

	ICartService interface {
		Create(req requests.CreateCart, roleres responses.BaseRole) (*responses.Cart, error)
		Read() ([]*responses.Cart, error)
		Update(cid uuid.UUID, req requests.UpdateCart) (*responses.Cart, error)
		Delete(id uuid.UUID) error

		GetByID(id uuid.UUID) (*responses.Cart, error)
		FindByRoleRes(roleres responses.BaseRole) ([]*responses.Cart, error)
	}
)

func NewCartService(
	rcart cartrepo.ICartRepository,
	rcare caregiverrepo.ICaregiverRepository,
	rmemb memberrepo.IMemberRepository,
	rmeal mealrepo.IMealRepository,
	sbsrl baseroleservice.IBaseRoleService,
) *CartService {
	return &CartService{
		rcart: rcart,
		rcare: rcare,
		rmemb: rmemb,
		rmeal: rmeal,

		sbsrl: sbsrl,
	}
}

func (s *CartService) Create(req requests.CreateCart, roleres responses.BaseRole) (*responses.Cart, error) {
	var (
		m   *models.Member
		err error
	)

	m, err = s.sbsrl.GetMemberByBaseRole(roleres)
	if err != nil {
		return nil, err
	}

	meal, err := s.rmeal.GetByID(req.MealID)
	if err != nil {
		return nil, err
	}

	// * convert request to model
	cart, err := req.ToModel(*m, *meal)
	if err != nil {
		utlogger.Error(err)
		return nil, consttypes.ErrConvertFailed
	}

	membres, err := m.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	// * try to get existing cart by meal ID and reference
	existingcart, err := s.rcart.GetByMealIDAndMemberID(cart.MemberID, cart.MealID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, consttypes.ErrGettingCart
	}

	// * if an existing cart is found, update its quantity
	if existingcart != nil {
		existingcart.Quantity += req.Quantity

		cart, err = s.rcart.Update(*existingcart)
		if err != nil {
			return nil, consttypes.ErrFailedToUpdateCart
		}
	} else {
		// * create a new cart
		cart, err = s.rcart.Create(*cart)
		if err != nil {
			return nil, consttypes.ErrFailedToCreateCart
		}
	}

	// * convert cart to response
	cartres, err := cart.ToResponse(membres)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return cartres, nil
}

func (s *CartService) Read() ([]*responses.Cart, error) {
	var (
		cartreses []*responses.Cart
	)

	carts, err := s.rcart.Read()
	if err != nil {
		return nil, consttypes.ErrFailedToReadCart
	}

	for _, cart := range carts {
		memb, err := s.rmemb.GetByID(cart.MemberID)
		if err != nil {
			return nil, consttypes.ErrCartNotFound
		}

		membres, err := memb.ToResponse()
		if err != nil {
			return nil, consttypes.ErrConvertFailed
		}

		cartres, err := cart.ToResponse(membres)
		if err != nil {
			return nil, consttypes.ErrConvertFailed
		}

		cartreses = append(cartreses, cartres)
	}

	return cartreses, nil
}

func (s *CartService) Update(cid uuid.UUID, req requests.UpdateCart) (*responses.Cart, error) {
	cart, err := s.rcart.GetByID(cid)
	if err != nil {
		return nil, consttypes.ErrCartNotFound
	}

	cart, err = req.ToModel(*cart)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	// * check if the result is 0 or not
	// * if it is 0, then delete the cart
	tempquantity := int(cart.Quantity) + req.Quantity
	if tempquantity < 0 || tempquantity == 0 {
		err := s.rcart.Delete(*cart)
		if err != nil {
			return nil, consttypes.ErrFailedToDeleteCart
		}

		// TODO: implement a correct way to return success deletion
		return nil, nil
	}

	// * if the value is not 0, then assign the variable
	cart.Quantity = tempquantity

	cart, err = s.rcart.Update(*cart)
	if err != nil {
		return nil, consttypes.ErrFailedToUpdateCart
	}

	memb, err := s.rmemb.GetByID(cart.MemberID)
	if err != nil {
		return nil, consttypes.ErrMemberNotFound
	}

	membres, err := memb.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	cartres, err := cart.ToResponse(membres)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return cartres, nil
}

func (s *CartService) Delete(id uuid.UUID) error {
	cart, err := s.rcart.GetByID(id)
	if err != nil {
		return consttypes.ErrCartNotFound
	}

	err = s.rcart.Delete(*cart)
	if err != nil {
		return consttypes.ErrFailedToDeleteCart
	}

	return nil
}

func (s *CartService) GetByID(id uuid.UUID) (*responses.Cart, error) {
	cart, err := s.rcart.GetByID(id)
	if err != nil {
		return nil, consttypes.ErrCartNotFound
	}

	memb, err := s.rmemb.GetByID(cart.MemberID)
	if err != nil {
		return nil, consttypes.ErrMemberNotFound
	}

	membres, err := memb.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	cartres, err := cart.ToResponse(membres)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return cartres, nil
}

func (s *CartService) FindByRoleRes(roleres responses.BaseRole) ([]*responses.Cart, error) {
	var (
		cartreses []*responses.Cart
	)

	m, err := s.sbsrl.GetMemberByBaseRole(roleres)
	if err != nil {
		return nil, consttypes.ErrMemberNotFound
	}

	mres, err := m.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	carts, err := s.rcart.FindByMemberID(m.ID)
	if err != nil {
		return nil, consttypes.ErrGettingCart
	}

	for _, cart := range carts {
		cartres, err := cart.ToResponse(mres)
		if err != nil {
			return nil, err
		}

		cartreses = append(cartreses, cartres)
	}

	return cartreses, nil
}
