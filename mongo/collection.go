package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/renkman/mongotui/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collection struct {
	currentCollection *mongo.Collection
}

var Collection *collection = &collection{}

func (collection *collection) SetCollection(name string) {
	collection.currentCollection = Database.currentDatabase.Collection(name)
}

func (collection *collection) Find(ctx context.Context, filter []byte, sort []byte, project []byte, limit int64, skip int64) chan models.QueryResult {
	ch := make(chan models.QueryResult)

	go func() {
		if collection.currentCollection == nil {
			ch <- models.QueryResult{nil, fmt.Errorf("No collection selected"), time.Since(time.Now())}
			return
		}

		if filter == nil || len(filter) == 0 {
			filter = []byte(`{}`)
		}

		if sort == nil || len(sort) == 0 {
			sort = []byte(`{}`)
		}

		if project == nil || len(project) == 0 {
			project = []byte(`{}`)
		}

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		filterBson, err := unmarshal(filter)
		if err != nil {
			ch <- models.QueryResult{nil, err, time.Since(time.Now())}
			return
		}

		sortBson, err := unmarshal(sort)
		if err != nil {
			ch <- models.QueryResult{nil, err, time.Since(time.Now())}
			return
		}

		projectBson, err := unmarshal(project)
		if err != nil {
			ch <- models.QueryResult{nil, err, time.Since(time.Now())}
			return
		}

		opts := options.Find().SetSort(sortBson).
			SetProjection(projectBson).
			SetLimit(limit).
			SetSkip(skip)

		start := time.Now()
		cursor, err := collection.currentCollection.Find(ctx, filterBson, opts)
		if err != nil {
			ch <- models.QueryResult{nil, err, time.Since(time.Now())}
			return
		}

		var result []map[string]interface{}

		cursor.RemainingBatchLength()
		err = cursor.All(ctx, &result)
		stop := time.Now()
		elapsed := stop.Sub(start)
		if err != nil {
			ch <- models.QueryResult{nil, err, elapsed}
			return
		}

		ch <- models.QueryResult{result, err, elapsed}
	}()
	return ch
}

func (collection *collection) Count(ctx context.Context, filter []byte) chan models.CountResult {
	ch := make(chan models.CountResult)

	go func() {
		if collection.currentCollection == nil {
			ch <- models.CountResult{0, fmt.Errorf("No collection selected"), time.Since(time.Now())}
			return
		}

		if filter == nil || len(filter) == 0 {
			start := time.Now()
			count, err := collection.currentCollection.EstimatedDocumentCount(ctx)
			stop := time.Now()
			duration := stop.Sub(start)

			if err != nil {
				ch <- models.CountResult{0, err, duration}
				return
			}
			ch <- models.CountResult{count, err, duration}
			return
		}

		filterBson, err := unmarshal(filter)
		if err != nil {
			ch <- models.CountResult{0, err, time.Since(time.Now())}
		}

		start := time.Now()
		count, err := collection.currentCollection.CountDocuments(ctx, filterBson)
		stop := time.Now()
		duration := stop.Sub(start)

		if err != nil {
			ch <- models.CountResult{0, err, duration}
			return
		}
		ch <- models.CountResult{count, err, duration}
	}()
	return ch
}

func (collection *collection) EstimatedCount(ctx context.Context) chan models.CountResult {
	ch := make(chan models.CountResult)

	go func() {
		if collection.currentCollection == nil {
			ch <- models.CountResult{0, fmt.Errorf("No collection selected"), time.Since(time.Now())}
			return
		}

		start := time.Now()
		count, err := collection.currentCollection.EstimatedDocumentCount(ctx)
		stop := time.Now()
		duration := stop.Sub(start)

		if err != nil {
			ch <- models.CountResult{0, err, duration}
			return
		}
		ch <- models.CountResult{count, err, duration}
	}()
	return ch
}

func unmarshal(command []byte) (interface{}, error) {
	if command == nil {
		return nil, nil
	}
	var commandBson interface{}
	err := bson.UnmarshalExtJSON(command, true, &commandBson)
	if err != nil {
		return nil, err
	}
	return commandBson, nil
}
