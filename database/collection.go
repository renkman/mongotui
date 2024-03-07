package database

import "context"

type Collection interface {
	SetCollection(name string)
	Find(ctx context.Context, filter []byte, sort []byte, project []byte) ([]map[string]interface{}, error)
}
