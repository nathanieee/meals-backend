package utpagination

import "time"

type (
	Pagination struct {
		Limit      int         `json:"limit,omitempty"`
		Page       int         `json:"page,omitempty"`
		Direction  string      `json:"direction,omitempty"`
		Sort       string      `json:"sort,omitempty"`
		Search     string      `json:"search,omitempty"`
		Filter     Filter      `json:"-"`
		TotalDatas int64       `json:"total_datas"`
		TotalPages int         `json:"total_pages"`
		Data       interface{} `json:"datas"`
	}

	Filter struct {
		CreatedFrom time.Time
		CreatedTo   time.Time
	}
)
