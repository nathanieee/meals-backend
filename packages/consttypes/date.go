package consttypes

import "time"

var (
	DateNow = time.Now().UTC()
)

const (
	DATEFORMAT string = "2006-01-02"
)
