package consttypes

import (
	"os"
	"time"
)

var (
	TIMEZONE, _ = time.LoadLocation(os.Getenv("API_TIMEZONE"))
)

const (
	DATEFORMAT                string = "2006-01-02"
	DATETIMEHOURMINUTESFORMAT string = "2006-01-02 15:04"
)

func TimeNow() time.Time {
	return time.Now().In(TIMEZONE)
}
