package memberallergyrepo

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	SELECTED_FIELDS = `
		id,
		member_id,
		allergy_id,
		created_at,
		updated_at
	`
)

type (
	MemberAllergyRepository struct {
		db *gorm.DB
	}

	IMemberAllergyRepository interface {
		Create(mall models.MemberAllergy) (*models.MemberAllergy, error)
		Delete(mall models.MemberAllergy) error
		GetByID(id uuid.UUID) (*models.MemberAllergy, error)
		GetByMemberIDAndAllergyID(mid uuid.UUID, aid uuid.UUID) (*models.MemberAllergy, error)
		GetByMemberID(mid uuid.UUID) ([]*models.MemberAllergy, error)
	}
)

func NewMemberAllergyRepository(db *gorm.DB) *MemberAllergyRepository {
	return &MemberAllergyRepository{db: db}
}

func (r *MemberAllergyRepository) omit() *gorm.DB {
	return r.db.
		Omit(
			"Allergy",
		)
}

func (r *MemberAllergyRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations)
}

func (r *MemberAllergyRepository) Create(mall models.MemberAllergy) (*models.MemberAllergy, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&mall).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	mallnew, err := r.GetByID(mall.ID)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mallnew, nil
}

func (r *MemberAllergyRepository) Delete(mall models.MemberAllergy) error {
	err := r.
		omit().
		Delete(&mall).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *MemberAllergyRepository) GetByID(id uuid.UUID) (*models.MemberAllergy, error) {
	var (
		mall *models.MemberAllergy
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("id = ?", id).
		First(&mall).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mall, nil
}

func (r *MemberAllergyRepository) GetByMemberIDAndAllergyID(mid uuid.UUID, iid uuid.UUID) (*models.MemberAllergy, error) {
	var (
		mall *models.MemberAllergy
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("member_id = ? AND allergy_id = ?", mid, iid).
		First(&mall).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return nil, nil
}

func (r *MemberAllergyRepository) GetByMemberID(mid uuid.UUID) ([]*models.MemberAllergy, error) {
	var (
		malls []*models.MemberAllergy
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("member_id = ?", mid).
		Find(&malls).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return malls, nil
}
