package database

import "context"

type Collection interface {
	SetCollection(name string)
	Find(ctx context.Context, filter []byte, sort []byte, project []byte) ([]map[string]interface{}, error)
	Count(ctx context.Context, filter []byte) (int64, error)
	EstimatedCount(ctx context.Context) (int64, error)
}
