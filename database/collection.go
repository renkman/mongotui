package database

import (
	"context"

	"github.com/renkman/mongotui/models"
)

type Collection interface {
	SetCollection(name string)
	Find(ctx context.Context, filter []byte, sort []byte, project []byte) chan models.QueryResult
	Count(ctx context.Context, filter []byte) (int64, error)
	EstimatedCount(ctx context.Context) (int64, error)
}
