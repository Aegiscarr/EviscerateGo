package time

import "time"

func Str2Utc(s string) (int64, error) {
	// parse twitter time string into UTC seconds, unix-style
	// golang's layout string for the date/time format
	t, err := time.Parse("01/02 03:04:05PM '06 -0700", s)
	return t.Unix(), err
}
