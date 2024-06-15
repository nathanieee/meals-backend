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
	"project-skbackend/packages/utils/utrole"

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
		Create(req requests.CreateCart, roleres responses.BaseRole) (*responses.Cart, error)
		Read() ([]*responses.Cart, error)
		Update(cid uuid.UUID, req requests.UpdateCart) (*responses.Cart, error)
		Delete(id uuid.UUID) error

		GetByID(id uuid.UUID) (*responses.Cart, error)
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

func (s *CartService) Create(req requests.CreateCart, roleres responses.BaseRole) (*responses.Cart, error) {
	var (
		m   *models.Member
		err error
	)

	rid, rtype, ok := utrole.CartRoleCheck(roleres)
	if !ok {
		return nil, consttypes.ErrUserInvalidRole
	}

	if rtype == consttypes.UR_CAREGIVER {
		m, err = s.rmemb.GetByCaregiverID(rid)
		if err != nil {
			return nil, err
		}
	} else if rtype == consttypes.UR_MEMBER {
		m, err = s.rmemb.GetByID(rid)
		if err != nil {
			return nil, err
		}
	}

	// Convert request to model
	cart, err := req.ToModel(*m)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	memb, err := s.rmemb.GetByID(cart.MemberID)
	if err != nil {
		return nil, err
	}

	membres, err := memb.ToResponse()
	if err != nil {
		return nil, err
	}

	// Try to get existing cart by meal ID and reference
	existingcart, err := s.rcart.GetByMealIDAndMemberID(cart.MemberID, cart.MealID)
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
	cartres, err := cart.ToResponse(membres)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cartres, nil
}

func (s *CartService) Read() ([]*responses.Cart, error) {
	var (
		cartreses []*responses.Cart
	)

	carts, err := s.rcart.Read()
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	for _, cart := range carts {
		memb, err := s.rmemb.GetByID(cart.MemberID)
		if err != nil {
			return nil, err
		}

		membres, err := memb.ToResponse()
		if err != nil {
			return nil, err
		}

		cartres, err := cart.ToResponse(membres)
		if err != nil {
			utlogger.Error(err)
			return nil, err
		}

		cartreses = append(cartreses, cartres)
	}

	return cartreses, nil
}

func (s *CartService) Update(cid uuid.UUID, req requests.UpdateCart) (*responses.Cart, error) {
	cart, err := s.rcart.GetByID(cid)

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

	memb, err := s.rmemb.GetByID(cart.MemberID)
	if err != nil {
		return nil, err
	}

	membres, err := memb.ToResponse()
	if err != nil {
		return nil, err
	}

	cartres, err := cart.ToResponse(membres)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cartres, nil
}

func (s *CartService) Delete(id uuid.UUID) error {
	cart, err := s.rcart.GetByID(id)
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

func (s *CartService) GetByID(id uuid.UUID) (*responses.Cart, error) {
	cart, err := s.rcart.GetByID(id)
	if err != nil {
		return nil, err
	}

	memb, err := s.rmemb.GetByID(cart.MemberID)
	if err != nil {
		return nil, err
	}

	membres, err := memb.ToResponse()
	if err != nil {
		return nil, err
	}

	cartres, err := cart.ToResponse(membres)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cartres, nil
}
