package consttypes

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	DATEFORMAT                string = "2006-01-02"
	DATETIMEHOURMINUTESFORMAT string = "2006-01-02 15:04"
)

func TimeNow() time.Time {
	godotenv.Load()

	tz, _ := time.LoadLocation(os.Getenv("TZ"))

	return time.Now().In(tz)
}
