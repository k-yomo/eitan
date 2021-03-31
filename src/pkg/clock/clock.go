package clock

import (
	"time"

	"github.com/pkg/errors"
)

var Now = func() time.Time {
	return time.Now()
}

func SetTimeZone(timeZone string) error {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return errors.Wrap(err, "load timezone")
	}
	time.Local = loc
	return nil
}

func BeginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}
