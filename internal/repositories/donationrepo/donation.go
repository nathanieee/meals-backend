package donationrepo

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
		patron_id,
		value,
		status,
		created_at,
		updated_at
	`
)

type (
	DonationRepository struct {
		db *gorm.DB
	}

	IDonationRepository interface {
		Create(d models.Donation) (*models.Donation, error)
		Read() ([]*models.Donation, error)
		Update(d models.Donation) (*models.Donation, error)
		Delete(d models.Donation) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*models.Donation, error)
		FindByPatronID(pid uuid.UUID) ([]*models.Donation, error)
	}
)

func NewDonationRepository(db *gorm.DB) *DonationRepository {
	return &DonationRepository{db: db}
}

func (r *DonationRepository) omit() *gorm.DB {
	return r.db.Omit(
		"",
	)
}

func (r *DonationRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations)
}

func (r *DonationRepository) Create(d models.Donation) (*models.Donation, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&d).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	dnew, err := r.GetByID(d.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return dnew, nil
}

func (r *DonationRepository) Read() ([]*models.Donation, error) {
	var (
		d []*models.Donation
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&d).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return d, nil
}

func (r *DonationRepository) Update(d models.Donation) (*models.Donation, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&d).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	dnew, err := r.GetByID(d.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return dnew, nil
}

func (r *DonationRepository) Delete(d models.Donation) error {
	err := r.db.
		Delete(&d).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *DonationRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		d    []models.Donation
		dres []responses.Donation
	)

	result := r.
		preload().
		Model(&d).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(r.db.
				Where("name LIKE ?", p.Search).
				Or("description LIKE ?", p.Search),
			)
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.
			Where("date(created_at) between ? and ?",
				p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT),
				p.Filter.CreatedTo.Format(consttypes.DATEFORMAT),
			)
	}

	if p.Filter.Patron.ID != nil {
		result = result.
			Where(&models.Donation{PatronID: *p.Filter.Patron.ID})
	}

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&d, &p, result)).
		Find(&d)

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from the model to response
	copier.CopyWithOption(&dres, &d, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = dres
	return &p, nil
}

func (r *DonationRepository) GetByID(id uuid.UUID) (*models.Donation, error) {
	var (
		d *models.Donation
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("id = ?", id).
		First(&d).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return d, err
}

func (r *DonationRepository) FindByPatronID(pid uuid.UUID) ([]*models.Donation, error) {
	var (
		d []*models.Donation
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("patron_id = ?", pid).
		Find(&d).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return d, nil
}
