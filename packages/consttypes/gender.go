package consttypes

import "database/sql/driver"

type (
	Gender string
)

const (
	G_MALE   Gender = "Male"
	G_FEMALE Gender = "Female"
	G_OTHER  Gender = "Other"
)

func (self *Gender) Scan(value interface{}) error {
	*self = Gender(value.([]byte))
	return nil
}

func (self Gender) Value() (driver.Value, error) {
	return string(self), nil
}
