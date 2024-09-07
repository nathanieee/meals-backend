package partnerservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/ordermealrepo"
	"project-skbackend/internal/repositories/orderrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	PartnerService struct {
		rpart partnerrepo.IPartnerRepository
		rordr orderrepo.IOrderRepository
		ruser userrepo.IUserRepository
		rorme ordermealrepo.IOrderMealRepository
	}

	IPartnerService interface {
		Create(req requests.CreatePartner) (*responses.Partner, error)
		Read() ([]*models.Partner, error)
		Update(id uuid.UUID, req requests.UpdatePartner) (*responses.Partner, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*responses.Partner, error)

		// * order related
		FindOwnOrders(uid uuid.UUID, preq utpagination.Pagination) (*utpagination.Pagination, error)
		OrderConfirmed(oid uuid.UUID, uid uuid.UUID) error
		OrderBeingPrepared(oid uuid.UUID, uid uuid.UUID) error
		OrderPrepared(oid uuid.UUID, uid uuid.UUID) error
		OrderPickedUp(oid uuid.UUID, uid uuid.UUID) error
	}
)

func NewPartnerService(
	rpart partnerrepo.IPartnerRepository,
	rordr orderrepo.IOrderRepository,
	rorme ordermealrepo.IOrderMealRepository,
) *PartnerService {
	return &PartnerService{
		rpart: rpart,
		rordr: rordr,
		rorme: rorme,
	}
}

func (s *PartnerService) Create(req requests.CreatePartner) (*responses.Partner, error) {
	user, err := req.User.ToModel(consttypes.UR_PARTNER)
	if err != nil {
		return nil, err
	}

	partner, err := req.ToModel(*user)
	if err != nil {
		return nil, err
	}

	partner, err = s.rpart.Create(*partner)
	if err != nil {
		return nil, err
	}

	pres, err := partner.ToResponse()
	if err != nil {
		return nil, err
	}

	return pres, nil
}

func (s *PartnerService) Read() ([]*models.Partner, error) {
	partners, err := s.rpart.Read()
	if err != nil {
		return nil, err
	}

	return partners, nil
}

func (s *PartnerService) Update(id uuid.UUID, req requests.UpdatePartner) (*responses.Partner, error) {
	partner, err := s.rpart.GetByID(id)
	if err != nil {
		return nil, err
	}

	user, err := req.User.ToModel(partner.User, consttypes.UR_PARTNER)
	if err != nil {
		return nil, err
	}

	partner, err = req.ToModel(*partner, *user)
	if err != nil {
		return nil, err
	}

	partner, err = s.rpart.Update(*partner)
	if err != nil {
		return nil, err
	}

	pres, err := partner.ToResponse()
	if err != nil {
		return nil, err
	}

	return pres, nil
}

func (s *PartnerService) Delete(id uuid.UUID) error {
	partner, err := s.rpart.GetByID(id)
	if err != nil {
		return err
	}

	return s.rpart.Delete(*partner)
}

func (s *PartnerService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	partners, err := s.rpart.FindAll(preq)
	if err != nil {
		return nil, err
	}

	return partners, nil
}

func (s *PartnerService) GetByID(id uuid.UUID) (*responses.Partner, error) {
	partner, err := s.rpart.GetByID(id)
	if err != nil {
		return nil, err
	}

	pres, err := partner.ToResponse()
	if err != nil {
		return nil, err
	}

	return pres, nil
}

func (s *PartnerService) FindOwnOrders(uid uuid.UUID, preq utpagination.Pagination) (*utpagination.Pagination, error) {
	partner, err := s.rpart.GetByUserID(uid)
	if err != nil {
		return nil, err
	}

	// * assigning partner id to the filter
	preq.Filter.Partner.ID = &partner.ID

	ordermeals, err := s.rorme.FindAll(preq)
	if err != nil {
		return nil, err
	}

	return ordermeals, nil
}

func (s *PartnerService) OrderConfirmed(oid uuid.UUID, uid uuid.UUID) error {
	// * get the user who confirms the order
	user, err := s.ruser.GetByID(uid)
	if err != nil {
		return err
	}

	// * get the corresponding order
	order, err := s.rordr.GetByID(oid)
	if err != nil {
		return err
	}

	// * could only confirm the order if the order status is "placed"
	if order.Status != consttypes.OS_PLACED {
		return consttypes.ErrInvalidOrderStatus
	}

	// * update the order object to confirm the order
	conorder := order.OrderConfirmed(*user)

	// * update the order in the database
	_, err = s.rordr.Update(*conorder)
	if err != nil {
		return err
	}

	return nil
}

func (s *PartnerService) OrderBeingPrepared(oid uuid.UUID, uid uuid.UUID) error {
	// * get the user who confirms the order
	user, err := s.ruser.GetByID(uid)
	if err != nil {
		return err
	}

	// * get the corresponding order
	order, err := s.rordr.GetByID(oid)
	if err != nil {
		return err
	}

	// * could only confirm the order if the order status is "confirmed"
	if order.Status != consttypes.OS_CONFIRMED {
		return consttypes.ErrInvalidOrderStatus
	}

	// * update the order object to confirm the order
	beporder := order.OrderBeingPrepared(*user)

	// * update the order in the database
	_, err = s.rordr.Update(*beporder)
	if err != nil {
		return err
	}

	return nil
}

func (s *PartnerService) OrderPrepared(oid uuid.UUID, uid uuid.UUID) error {
	// * get the user who confirms the order
	user, err := s.ruser.GetByID(uid)
	if err != nil {
		return err
	}

	// * get the corresponding order
	order, err := s.rordr.GetByID(oid)
	if err != nil {
		return err
	}

	// * could only confirm the order if the order status is "being prepared"
	if order.Status != consttypes.OS_BEING_PREPARED {
		return consttypes.ErrInvalidOrderStatus
	}

	// * update the order object to confirm the order
	preporder := order.OrderPrepared(*user)

	// * update the order in the database
	_, err = s.rordr.Update(*preporder)
	if err != nil {
		return err
	}

	return nil
}

func (s *PartnerService) OrderPickedUp(oid uuid.UUID, uid uuid.UUID) error {
	// * get the user who confirms the order
	user, err := s.ruser.GetByID(uid)
	if err != nil {
		return err
	}

	// * get the corresponding order
	order, err := s.rordr.GetByID(oid)
	if err != nil {
		return err
	}

	// * could only confirm the order if the order status is "being prepared"
	if order.Status != consttypes.OS_PREPARED {
		return consttypes.ErrInvalidOrderStatus
	}

	// * update the order object to confirm the order
	preporder := order.OrderPickedUp(*user)

	// * update the order in the database
	_, err = s.rordr.Update(*preporder)
	if err != nil {
		return err
	}

	return nil
}
