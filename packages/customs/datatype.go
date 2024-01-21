package customs

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"project-skbackend/packages/consttypes"
	"strings"
	"time"
)

type (
	CDT_STRING string
	CDT_DATE   time.Time
)

/* ------------------------------- CDT_STRING ------------------------------- */
func (self *CDT_STRING) SuffixSpaceCheck() string {
	mstr := string(*self)

	if mstr != "" && !strings.HasSuffix(mstr, " ") {
		mstr += " "
	}

	return mstr
}

/* -------------------------------- CDT_DATE -------------------------------- */
func (self *CDT_DATE) UnmarshalJSON(data []byte) error {
	var datestr string
	if err := json.Unmarshal(data, &datestr); err != nil {
		return err
	}

	t, err := time.Parse(consttypes.DATEFORMAT, datestr)
	if err != nil {
		return err
	}

	*self = CDT_DATE(t)
	return nil
}

func (self CDT_DATE) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(self).Format(consttypes.DATEFORMAT))
}

func (self CDT_DATE) Value() (driver.Value, error) {
	return time.Time(self), nil
}

func (self *CDT_DATE) Scan(value interface{}) error {
	if value == nil {
		*self = CDT_DATE(time.Time{})
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("Unable to convert value to time.Time")
	}

	*self = CDT_DATE(t)
	return nil
}
