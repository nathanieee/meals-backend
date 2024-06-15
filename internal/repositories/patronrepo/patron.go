package patronrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/base"
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
		user_id,
		type,
		name,
		created_at,
		updated_at
	`
)

type (
	PatronRepository struct {
		db *gorm.DB
	}

	IPatronRepository interface {
		Create(p models.Patron) (*models.Patron, error)
		Read() ([]*models.Patron, error)
		Update(p models.Patron) (*models.Patron, error)
		Delete(p models.Patron) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*models.Patron, error)
		FindByUserID(id uuid.UUID) (*models.Patron, error)
	}
)

func NewPatronRepository(db *gorm.DB) *PatronRepository {
	return &PatronRepository{db: db}
}

func (r *PatronRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("User.Address").
		Preload("User.Image.Image").
		Preload("Donations")
}

func (r *PatronRepository) Create(p models.Patron) (*models.Patron, error) {
	err := r.db.
		Create(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	pnew, err := r.FindByID(p.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return pnew, nil
}

func (r *PatronRepository) Read() ([]*models.Patron, error) {
	var (
		p []*models.Patron
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return p, nil
}

func (r *PatronRepository) Update(p models.Patron) (*models.Patron, error) {
	err := r.db.
		Save(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	pnew, err := r.FindByID(p.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return pnew, nil
}

func (r *PatronRepository) Delete(p models.Patron) error {
	err := r.db.
		Delete(&p).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *PatronRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		ps     []models.Patron
		preses []responses.Patron
	)

	result := r.
		preload().
		Model(&ps).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(r.db.
				Where("name LIKE ?", p.Search),
			)
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
		Scopes(paginationrepo.Paginate(&ps, &p, result)).
		Find(&ps)

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&preses, &ps, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = preses
	return &p, nil
}

func (r *PatronRepository) FindByID(id uuid.UUID) (*models.Patron, error) {
	var (
		p *models.Patron
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Patron{Model: base.Model{ID: id}}).
		First(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return p, nil
}

func (r *PatronRepository) FindByUserID(id uuid.UUID) (*models.Patron, error) {
	var (
		p *models.Patron
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Patron{User: models.User{Model: base.Model{ID: id}}}).
		First(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return p, nil
}
