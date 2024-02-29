<<<<<<< HEAD
package main

import (
	"time"
)

func Str2utc(s string) (int64, error) {
	// parse twitter time string into UTC seconds, unix-style
	// golang's layout string for the date/time format
	t, err := time.Parse("01/02 03:04:05PM '06 -0700", s)
	return t.Unix(), err
}

func Utc2snowflake(stamp int64) int64 {
	return (stamp*1000 - 1288834974657) << 22
}

func Snowflake2utc(sf int64) float64 {
	return float64((sf>>22)+1288834974657) / 1000.0
}

func Str2utcms(s string) int64 {
	t, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", s)
	return t.UnixNano() / int64(time.Millisecond)
}

func Snowflake2utcms(sf int64) int64 {
	return (sf >> 22) + 1288834974657
}

// really is the best way to get utc timestamp?
//
//	(minus changing your box to be UTC)
func Utcnow() int64 {
	return time.Now().Unix()
}

// hey thanks chatgpt lmao
=======
package main

import (
	"time"
)

func Str2utc(s string) (int64, error) {
	// parse twitter time string into UTC seconds, unix-style
	// golang's layout string for the date/time format
	t, err := time.Parse("01/02 03:04:05PM '06 -0700", s)
	return t.Unix(), err
}

func Utc2snowflake(stamp int64) int64 {
	return (stamp*1000 - 1288834974657) << 22
}

func Snowflake2utc(sf int64) float64 {
	return float64((sf>>22)+1288834974657) / 1000.0
}

func Str2utcms(s string) int64 {
	t, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", s)
	return t.UnixNano() / int64(time.Millisecond)
}

func Snowflake2utcms(sf int64) int64 {
	return (sf >> 22) + 1288834974657
}

// really is the best way to get utc timestamp?
//
//	(minus changing your box to be UTC)
func Utcnow() int64 {
	return time.Now().Unix()
}

// hey thanks chatgpt lmao
>>>>>>> 42aba7d (the 'i fixed some timeouts' update)
