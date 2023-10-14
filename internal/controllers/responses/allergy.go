package responses

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

type (
	AllergyResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name" gorm:"not null" binding:"required" example:"Lactose Intolerant"`
	}
)

func (alres *AllergyResponse) IsEmpty() bool {
	return cmp.Equal(alres, AllergyResponse{})
}
