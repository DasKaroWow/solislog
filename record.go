package solislog

import (
	"time"
)

type record struct {
	time    time.Time
	level   Level
	message string
	extra   Extra
}
