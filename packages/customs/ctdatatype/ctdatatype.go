package ctdatatype

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

// ! -------------------------------- CDT_DATE -------------------------------- ! //
func (date *CDT_DATE) UnmarshalJSON(data []byte) error {
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

	*date = CDT_DATE{Time: t}
	return nil
}

func (date CDT_DATE) MarshalJSON() ([]byte, error) {
	return json.Marshal(date.Format(consttypes.DATEFORMAT))
}

func (date CDT_DATE) ToTime() (time.Time, error) {
	timedate := date

	if timedate.IsZero() {
		return time.Time{}, nil
	}

	return date.Time, nil
}

func (date CDT_DATE) Value() (driver.Value, error) {
	return date.Time, nil
}

func (date *CDT_DATE) Scan(value any) error {
	if value == nil {
		*date = CDT_DATE{}
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("unable to convert value to time.Time")
	}

	*date = CDT_DATE{Time: t}
	return nil
}
