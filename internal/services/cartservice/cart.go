package cartservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/cartrepo"
	"project-skbackend/internal/repositories/memberrepo"
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
	}

	ICartService interface {
		Create(req requests.CreateCart) (*responses.Cart, error)
		Read() ([]*responses.Cart, error)
		Update(req requests.UpdateCart) (*responses.Cart, error)
		Delete(id uuid.UUID) error

		FindByID(id uuid.UUID) (*responses.Cart, error)
		GetCartByMealIDAndReference(mid uuid.UUID, rid uuid.UUID, rtype consttypes.UserRole) (*responses.Cart, error)
		GetCartReferenceObject(cart models.Cart) (*responses.Member, *responses.Caregiver, error)
	}
)

func NewCartService(
	rcart cartrepo.ICartRepository,
	rcare caregiverrepo.ICaregiverRepository,
	rmemb memberrepo.IMemberRepository,
) *CartService {
	return &CartService{
		rcart: rcart,
		rcare: rcare,
		rmemb: rmemb,
	}
}

func (s *CartService) Create(req requests.CreateCart) (*responses.Cart, error) {
	// Convert request to model
	cart, err := req.ToModel()
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// Get cart reference objects
	membres, careres, err := s.GetCartReferenceObject(*cart)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// Try to get existing cart by meal ID and reference
	existingcart, err := s.rcart.GetCartByMealIDAndReference(cart.MealID, cart.ReferenceID, cart.ReferenceType)
	if err != nil && err != gorm.ErrRecordNotFound {
		utlogger.Error(err)
		return nil, err
	}

	// If an existing cart is found, update its quantity
	if existingcart != nil {
		existingcart.Quantity += req.Quantity

		cart, err = s.rcart.Update(*existingcart)
		if err != nil {
			utlogger.Error(err)
			return nil, err
		}
	} else {
		// Create a new cart
		cart, err = s.rcart.Create(*cart)
		if err != nil {
			utlogger.Error(err)
			return nil, err
		}
	}

	// Convert cart to response
	cartres, err := cart.ToResponse(membres, careres)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cartres, nil
}

func (s *CartService) Read() ([]*responses.Cart, error) {
	carts, err := s.rcart.Read()
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cartreses := make([]*responses.Cart, 0, len(carts))
	for _, cart := range carts {
		membres, careres, err := s.GetCartReferenceObject(*cart)
		if err != nil {
			utlogger.Error(err)
			return nil, err
		}

		cartres, err := cart.ToResponse(membres, careres)
		if err != nil {
			utlogger.Error(err)
			return nil, err
		}

		cartreses = append(cartreses, cartres)
	}

	return cartreses, nil
}

func (s *CartService) Update(req requests.UpdateCart) (*responses.Cart, error) {
	cart, err := s.rcart.GetCartByMealIDAndReference(req.MealID, req.ReferenceID, req.ReferenceType)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cart, err = req.ToModel(*cart)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cart.Quantity = cart.Quantity + req.Quantity

	cart, err = s.rcart.Update(*cart)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	membres, careres, err := s.GetCartReferenceObject(*cart)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cartres, err := cart.ToResponse(membres, careres)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cartres, nil
}

func (s *CartService) Delete(id uuid.UUID) error {
	cart, err := s.rcart.FindByID(id)
	if err != nil {
		utlogger.Error(err)
		return err
	}

	err = s.rcart.Delete(*cart)
	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (s *CartService) FindByID(id uuid.UUID) (*responses.Cart, error) {
	cart, err := s.rcart.FindByID(id)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	membres, careres, err := s.GetCartReferenceObject(*cart)
	if err != nil {
		return nil, err
	}

	cartres, err := cart.ToResponse(membres, careres)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cartres, nil
}

func (s *CartService) GetCartReferenceObject(cart models.Cart) (*responses.Member, *responses.Caregiver, error) {
	return s.rcart.GetCartReferenceObject(cart)
}

func (s *CartService) GetCartByMealIDAndReference(mid uuid.UUID, rid uuid.UUID, rtype consttypes.UserRole) (*responses.Cart, error) {
	cart, err := s.rcart.GetCartByMealIDAndReference(mid, rid, rtype)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	membres, careres, err := s.GetCartReferenceObject(*cart)
	if err != nil {
		return nil, err
	}

	cartres, err := cart.ToResponse(membres, careres)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cartres, nil
}
