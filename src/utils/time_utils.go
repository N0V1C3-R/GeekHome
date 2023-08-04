package utils

import (
	"github.com/araddon/dateparse"
	"time"
)

func GetCurrentTime() time.Time {
	return time.Now()
}

func ConvertToUnixTime(t time.Time) int64 {
	return t.Unix()
}

func ConvertToMilliTime(t time.Time) int64 {
	return t.UnixMilli()
}

func ConvertToMicroTime(t time.Time) int64 {
	return t.UnixMicro()
}

func ConvertToNanoTime(t time.Time) int64 {
	return t.UnixNano()
}

func FormatCurrentTime(t time.Time) string {
	return t.Format("2006-01-02")
}

func ParseTimeString(t string) *time.Time {
	tTime, err := dateparse.ParseAny(t)
	if err != nil {
		return nil
	}
	return &tTime
}
