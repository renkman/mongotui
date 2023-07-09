package database

import "context"

type Collection interface {
	Find(ctx context.Context, filter []byte, sort []byte, project []byte) ([]map[string]interface{}, error)
}
