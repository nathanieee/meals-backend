package utpagination

import (
	"time"

	"github.com/google/uuid"
)

type (
	Pagination struct {
		Limit      int    `json:"limit,omitempty"`
		Page       int    `json:"page,omitempty"`
		Direction  string `json:"direction,omitempty"`
		Sort       string `json:"sort,omitempty"`
		Search     string `json:"search,omitempty"`
		Filter     Filter `json:"-"`
		TotalDatas int64  `json:"total_datas"`
		TotalPages int    `json:"total_pages"`
		Data       any    `json:"datas"`
	}

	Filter struct {
		CreatedFrom time.Time
		CreatedTo   time.Time
		Meal        Meal
		Partner     Partner
		Patron      Patron
	}

	Patron struct {
		ID *uuid.UUID `json:"patron_id"`
	}

	Meal struct {
		ID *uuid.UUID `json:"meal_id"`
	}

	Partner struct {
		ID *uuid.UUID `json:"partner_id"`
	}
)
