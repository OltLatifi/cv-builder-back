package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type CustomDate struct {
	time.Time
}

const dateLayout = "2006-01-02"

func (d *CustomDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	date, err := time.Parse(`"`+dateLayout+`"`, s)
	if err != nil {
		return err
	}
	d.Time = date
	return nil
}

func (d CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.Time.Format(dateLayout))), nil
}

// Conversion methods for GORM:
func (d CustomDate) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d *CustomDate) Scan(value interface{}) error {
	d.Time = value.(time.Time)
	return nil
}
