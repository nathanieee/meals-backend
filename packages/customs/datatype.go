package customs

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"project-skbackend/packages/consttypes"
	"time"
)

type (
	CDT_DATE struct{ time.Time }
)

/* -------------------------------- CDT_DATE -------------------------------- */
func (self *CDT_DATE) UnmarshalJSON(data []byte) error {
	var (
		datestr string
	)

	if err := json.Unmarshal(data, &datestr); err != nil {
		return err
	}

	t, err := time.Parse(consttypes.DATEFORMAT, datestr)
	if err != nil {
		return err
	}

	*self = CDT_DATE{Time: t}
	return nil
}

func (self CDT_DATE) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Format(consttypes.DATEFORMAT))
}

func (self CDT_DATE) ToTime() (time.Time, error) {
	timeself := self

	if timeself.IsZero() {
		return time.Time{}, nil
	}

	return self.Time, nil
}

func (self CDT_DATE) Value() (driver.Value, error) {
	return self.Time, nil
}

func (self *CDT_DATE) Scan(value any) error {
	if value == nil {
		*self = CDT_DATE{}
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("unable to convert value to time.Time")
	}

	*self = CDT_DATE{Time: t}
	return nil
}
