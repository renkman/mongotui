package models

import "time"

type QueryResult struct {
	Result   []map[string]interface{}
	Error    error
	Duration time.Duration
}
