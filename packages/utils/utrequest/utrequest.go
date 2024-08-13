package utrequest

import (
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utpagination"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobeam/stringy"
)

func CheckWhitelistUrl(url string) bool {
	splittedUrl := strings.Split(url, "api/v1/")
	whitelistedUrl := map[string]bool{
		"auth/refresh-token": true,
	}

	return whitelistedUrl[splittedUrl[1]]
}

func GeneratePaginationFromRequest(ctx *gin.Context) utpagination.Pagination {
	var (
		limit       = 10
		page        = 1
		search      = ""
		sort        = "id"
		direction   = "asc"
		createdFrom time.Time
		createdTo   time.Time
	)

	query := ctx.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		case "search":
			search = queryValue
		case "sort":
			if queryValue != "" {
				str := stringy.New(queryValue)
				snakeStr := str.SnakeCase("?", "")
				sort = snakeStr.ToLower()
			}
		case "direction":
			queryValue = strings.ToLower(queryValue)
			if queryValue == "asc" || queryValue == "desc" {
				direction = queryValue
			}
		case "created-from":
			queryValue = strings.ToLower(queryValue)
			if queryValue != "" {
				date, err := time.Parse(consttypes.DATEFORMAT, queryValue)
				if err == nil {
					createdFrom = date
				}
			}
		case "created-to":
			queryValue = strings.ToLower(queryValue)
			if queryValue != "" {
				date, err := time.Parse(consttypes.DATEFORMAT, queryValue)
				if err == nil {
					createdTo = date
				}
			}
		}
	}

	return utpagination.Pagination{
		Limit:     limit,
		Page:      page,
		Sort:      sort,
		Direction: direction,
		Search:    search,
		Filter: utpagination.Filter{
			CreatedFrom: createdFrom,
			CreatedTo:   createdTo,
		},
	}
}
