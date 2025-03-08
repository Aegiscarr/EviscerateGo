package time

func Snowflake2Utc(sf int64) float64 {
	return float64((sf>>22)+1288834974657) / 1000.0
}
