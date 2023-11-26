package pcsm

import (
	"strings"
	"time"
)

// UnmarshalJSON handles reading the C# formatted timestamp and
// converting it to a more standard format.
func (t *Time) UnmarshalJSON(d []byte) error {
	// 2023-11-25T22:35:54.3880495+00:00
	tp, err := time.Parse("2006-01-02T15:04:05.9999Z07:00", strings.Trim(string(d), "\""))
	t.Time = tp
	return err
}

// ToTime returns the underlying type object
func (t Time) ToTime() time.Time {
	return t.Time
}
