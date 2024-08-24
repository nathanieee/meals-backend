package donationproofrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
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
		donation_id,
		image_id,
		created_at,
		updated_at
	`
)

type (
	DonationProofRepository struct {
		db *gorm.DB
	}

	IDonationProofRepository interface {
		Create(dp models.DonationProof) (*models.DonationProof, error)
		Read() ([]*models.DonationProof, error)
		Update(dp models.DonationProof) (*models.DonationProof, error)
		Delete(dp models.DonationProof) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(dpid uuid.UUID) (*models.DonationProof, error)
		FindByDonationID(did uuid.UUID) ([]*models.DonationProof, error)
	}
)

func NewDonationProofRepository(db *gorm.DB) *DonationProofRepository {
	return &DonationProofRepository{db: db}
}

func (r *DonationProofRepository) omit() *gorm.DB {
	return r.db.Omit(
		"",
	)
}

func (r *DonationProofRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations)
}

func (r *DonationProofRepository) Create(dp models.DonationProof) (*models.DonationProof, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&dp).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	dpnew, err := r.GetByID(dp.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return dpnew, nil
}

func (r *DonationProofRepository) Read() ([]*models.DonationProof, error) {
	var (
		dp []*models.DonationProof
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&dp).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return dp, nil
}

func (r *DonationProofRepository) Update(dp models.DonationProof) (*models.DonationProof, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&dp).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	dpnew, err := r.GetByID(dp.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return dpnew, nil
}

func (r *DonationProofRepository) Delete(dp models.DonationProof) error {
	err := r.db.
		Delete(&dp).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *DonationProofRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		dp    []models.DonationProof
		dpres []responses.DonationProof
	)

	result := r.
		preload().
		Model(&dp).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		// TODO: add a like query here
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.
			Where("date(created_at) between ? and ?",
				p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT),
				p.Filter.CreatedTo.Format(consttypes.DATEFORMAT),
			)
	}

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&dp, &p, result)).
		Find(&dp)

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from the model to response
	copier.CopyWithOption(&dpres, &dp, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = dpres
	return &p, nil
}

func (r *DonationProofRepository) GetByID(dpid uuid.UUID) (*models.DonationProof, error) {
	var (
		dp *models.DonationProof
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("id = ?", dpid).
		First(&dp).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return dp, nil
}

func (r *DonationProofRepository) FindByDonationID(did uuid.UUID) ([]*models.DonationProof, error) {
	var (
		dp []*models.DonationProof
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("donation_id = ?", did).
		Find(&dp).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return dp, nil
}
