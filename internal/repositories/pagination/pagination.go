package pagination

import (
	"fmt"
	"math"
	"project-skbackend/packages/utils"

	"gorm.io/gorm"
)

func getLimit(p *utils.Pagination) int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func getPage(p *utils.Pagination) int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func getSort(p *utils.Pagination) string {
	builtQuery := "id desc"
	direction := "asc"

	if p.Direction != "" {
		direction = p.Direction
	}

	if p.Sort != "" {
		builtQuery = fmt.Sprintf("%s %s", p.Sort, direction)
	}

	return builtQuery
}

func getOffset(p *utils.Pagination) int {
	return (getPage(p) - 1) * getLimit(p)
}

func Paginate(utils interface{}, pagination *utils.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalDatas int64
	db.Count(&totalDatas)

	totalPages := int(math.Ceil(float64(totalDatas) / float64(pagination.Limit)))

	pagination.TotalDatas = totalDatas
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(getOffset(pagination)).Limit(getLimit(pagination)).Order(getSort(pagination))
	}
}
