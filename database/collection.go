package database

import (
	"context"

	"github.com/renkman/mongotui/models"
)

type Collection interface {
	SetCollection(name string)
	Find(ctx context.Context, filter []byte, sort []byte, project []byte, limit int64, skip int64) chan models.QueryResult
	Count(ctx context.Context, filter []byte) chan models.CountResult
	EstimatedCount(ctx context.Context) chan models.CountResult
}
