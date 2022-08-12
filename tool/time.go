package tool

import "time"

const (
	DefaultTimeFormatter = "2006-01-02:15:04:05"
	DateFormatter        = "20060102"
	DatePlusFormatter    = "2006-01-02"
	TimeFormatter        = "150405"
	TimePlusFormatter    = "15:04:05"
)

func FormatTime(t time.Time) (s string) {
	s = t.Format(DefaultTimeFormatter)
	return
}
