package memberillnessrepo

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
		illness_id,
		created_at,
		updated_at
	`
)

type (
	MemberIllnessRepository struct {
		db *gorm.DB
	}

	IMemberIllnessRepository interface {
		Create(mill models.MemberIllness) (*models.MemberIllness, error)
		Delete(mill models.MemberIllness) error
		GetByID(id uuid.UUID) (*models.MemberIllness, error)
		GetByMemberIDAndIllnessID(mid uuid.UUID, iid uuid.UUID) (*models.MemberIllness, error)
		GetByMemberID(mid uuid.UUID) ([]*models.MemberIllness, error)
	}
)

func NewMemberIllnessRepository(db *gorm.DB) *MemberIllnessRepository {
	return &MemberIllnessRepository{db: db}
}

func (r *MemberIllnessRepository) omit() *gorm.DB {
	return r.db.
		Omit(
			"Illness",
		)
}

func (r *MemberIllnessRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations)
}

func (r *MemberIllnessRepository) Create(mill models.MemberIllness) (*models.MemberIllness, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&mill).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	millnew, err := r.GetByID(mill.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return millnew, nil
}

func (r *MemberIllnessRepository) Delete(mill models.MemberIllness) error {
	err := r.
		omit().
		Delete(&mill).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}
func (r *MemberIllnessRepository) GetByID(id uuid.UUID) (*models.MemberIllness, error) {
	var (
		mill *models.MemberIllness
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("id = ?", id).
		First(&mill).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mill, nil
}

func (r *MemberIllnessRepository) GetByMemberIDAndIllnessID(mid uuid.UUID, iid uuid.UUID) (*models.MemberIllness, error) {
	var (
		mill *models.MemberIllness
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("member_id = ? AND illness_id = ?", mid, iid).
		First(&mill).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mill, nil
}

func (r *MemberIllnessRepository) GetByMemberID(mid uuid.UUID) ([]*models.MemberIllness, error) {
	var (
		mills []*models.MemberIllness
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where("member_id = ?", mid).
		Find(&mills).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mills, nil
}
