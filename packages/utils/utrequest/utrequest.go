package utrequest

import (
	"fmt"
	"math"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utmath"
	"project-skbackend/packages/utils/utpagination"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobeam/stringy"
)

var (
	suffixes [5]string
)

func CheckWhitelistUrl(url string) bool {
	splittedUrl := strings.Split(url, "api/v1/")
	whitelistedUrl := map[string]bool{
		"users/me":           true,
		"auth/refresh-token": true,
	}

	return whitelistedUrl[splittedUrl[1]]
}

func GeneratePaginationFromRequest(ctx *gin.Context) utpagination.Pagination {
	limit := 10
	page := 1
	search := ""
	sort := "id"
	direction := "asc"
	var createdFrom time.Time
	var createdTo time.Time

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

// TODO - needs to use this function to validate file sizes in the future
func GetReadableFileSize(size float64, ext string) error {
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"
	maxImageFile := 0.6
	maxImageFileSuffix := "MB"
	maxVideoSize := 30
	maxVideoSizeSuffix := "MB"
	maxSoundSize := 10
	maxSoundSizeSuffix := "MB"

	base := math.Log(size) / math.Log(1024)
	getSize := utmath.Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	if getSuffix == "KB" {
		getSize = math.Ceil((getSize/1000)*100) / 100
		getSuffix = "MB"
	}

	if maxImageFileSuffix == getSuffix || maxVideoSizeSuffix == getSuffix {
		switch ext {
		case "image":
			if float64(getSize) > maxImageFile || getSuffix == "GB" || getSuffix == "TB" {
				return fmt.Errorf("image size is too big. Maximum size is %f %s", float64(maxImageFile), maxImageFileSuffix)
			}
			return nil
		case "video":
			if int(getSize) > maxVideoSize || getSuffix == "GB" || getSuffix == "TB" {
				return fmt.Errorf("video size is too big. Maximum size is %d %s", maxVideoSize, maxVideoSizeSuffix)
			}
			return nil
		case "sound":
			if int(getSize) > maxSoundSize || getSuffix == "GB" || getSuffix == "TB" {
				return fmt.Errorf("sound size is too big. Maximum size is %d %s", maxSoundSize, maxSoundSizeSuffix)
			}
			return nil
		}
	}
	return nil
}
