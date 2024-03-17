package patronrepo

import "gorm.io/gorm"

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
)
