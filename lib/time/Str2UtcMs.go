package time

import "time"

func Str2UtcMs(s string) int64 {
	t, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", s)
	return t.UnixNano() / int64(time.Millisecond)
}
