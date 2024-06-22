package consttypes

import "time"

const (
	DATEFORMAT string = "2006-01-02"
)

func TimeNow() time.Time {
	return time.Now().UTC()
}
