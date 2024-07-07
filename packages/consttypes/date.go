package consttypes

import (
	"os"
	"time"
)

const (
	DATEFORMAT                string = "2006-01-02"
	DATETIMEHOURMINUTESFORMAT string = "2006-01-02 15:04"
)

func TimeNow() time.Time {
	tz, _ := time.LoadLocation(os.Getenv("API_TIMEZONE"))

	return time.Now().In(tz)
}
