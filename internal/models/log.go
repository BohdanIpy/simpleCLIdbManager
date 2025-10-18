package models

import "time"

type Log struct {
	Id         int64
	LogTime    time.Time
	LogMessage string
}
