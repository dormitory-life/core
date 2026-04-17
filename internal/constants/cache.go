package constants

import "time"

const (
	CacheDormitoriesKey = "dormitories"
)

const (
	DefaultDormitoryListTTL = time.Minute * 30
	DefaultDormitoryTTL     = time.Minute * 5
)
