package time

import "time"

func UtcNow() int64 {
	return time.Now().Unix()
}
