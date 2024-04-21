package models

import (
	"time"
)

type CountResult struct {
	Count    int64
	Error    error
	Duration time.Duration
}
