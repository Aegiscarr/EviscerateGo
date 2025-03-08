package time

func Utc2Snowflake(stamp int64) int64 {
	return (stamp*1000 - 1288834974657) << 22
}
