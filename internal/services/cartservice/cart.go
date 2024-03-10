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
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/utstring"

	"github.com/google/uuid"
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
		Update(id uuid.UUID, req requests.UpdateCart) (*responses.Cart, error)
		Delete(id uuid.UUID) error
		FindByID(id uuid.UUID) (*responses.Cart, error)
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
	cart, err := req.ToModel()
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	membresp, careresp, err := s.GetCartReferenceObject(*cart)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	cart, err = s.rcart.Create(*cart)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	cartresp, err := cart.ToResponse(membresp, careresp)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return cartresp, nil
}

func (s *CartService) Read() ([]*responses.Cart, error) {
	carts, err := s.rcart.Read()
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	cartresps := make([]*responses.Cart, 0, len(carts))
	for _, cart := range carts {
		membresp, careresp, err := s.GetCartReferenceObject(*cart)
		if err != nil {
			utlogger.LogError(err)
			return nil, err
		}

		cartresp, err := cart.ToResponse(membresp, careresp)
		if err != nil {
			utlogger.LogError(err)
			return nil, err
		}

		cartresps = append(cartresps, cartresp)
	}

	return cartresps, nil
}

func (s *CartService) Update(id uuid.UUID, req requests.UpdateCart) (*responses.Cart, error) {
	cart, err := s.rcart.FindByID(id)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	cart, err = req.ToModel(*cart)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	cart, err = s.rcart.Update(*cart)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	membresp, careresp, err := s.GetCartReferenceObject(*cart)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	cartresp, err := cart.ToResponse(membresp, careresp)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return cartresp, nil
}

func (s *CartService) Delete(id uuid.UUID) error {
	cart, err := s.rcart.FindByID(id)
	if err != nil {
		utlogger.LogError(err)
		return err
	}

	err = s.rcart.Delete(*cart)
	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (s *CartService) FindByID(id uuid.UUID) (*responses.Cart, error) {
	cart, err := s.rcart.FindByID(id)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	membresp, careresp, err := s.GetCartReferenceObject(*cart)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	cartresp, err := cart.ToResponse(membresp, careresp)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return cartresp, nil
}

func (s *CartService) GetCartReferenceObject(cart models.Cart) (*responses.Member, *responses.Caregiver, error) {
	var careresp *responses.Caregiver
	var membresp *responses.Member

	utstring.PrintJSON(cart)

	switch cart.ReferenceType {
	case consttypes.UR_CAREGIVER:
		caregiver, err := s.rcare.FindByID(cart.ReferenceID)
		if err != nil {
			utlogger.LogError(err)
			return nil, nil, err
		}

		careresp = caregiver.ToResponse()
	case consttypes.UR_MEMBER:
		member, err := s.rmemb.FindByID(cart.ReferenceID)
		if err != nil {
			utlogger.LogError(err)
			return nil, nil, err
		}

		membresp = member.ToResponse()
	default:
		return nil, nil, utresponse.ErrInvalidReference
	}

	return membresp, careresp, nil
}
