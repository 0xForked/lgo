package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type genericMongoQueryRepository[T any] struct {
	collection *mongo.Collection
}

func (repository genericMongoQueryRepository[T]) CountData(
	ctx context.Context,
	pipeline mongo.Pipeline,
) int64 {
	pipeline = append(pipeline, bson.D{{Key: "$count", Value: "count"}})
	cursor, err := repository.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return int64(0)
	}
	defer func() { _ = cursor.Close(ctx) }()
	var result struct {
		Count int64 `bson:"count"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return int64(0)
		}
	}
	if err := cursor.Err(); err != nil {
		return int64(0)
	}
	return result.Count
}

func (repository genericMongoQueryRepository[T]) AggregateData(
	ctx context.Context,
	pipeline mongo.Pipeline,
) ([]*T, error) {
	cursor, err := repository.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()
	var items []*T
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (repository genericMongoQueryRepository[T]) FindData(
	ctx context.Context,
	filter bson.M,
) (*T, error) {
	var item *T
	if err := repository.collection.FindOne(
		ctx, filter,
	).Decode(&item); err != nil {
		return nil, err
	}
	return item, nil
}

func (repository genericMongoQueryRepository[T]) DistinctData(
	ctx context.Context,
	filter bson.M,
	field string,
) ([]interface{}, error) {
	opts := options.Distinct()
	uniqueTeams, err := repository.collection.Distinct(ctx, field, filter, opts)
	if err != nil {
		return nil, err
	}
	return uniqueTeams, nil
}

func NewGenericMongoQueryRepository[T any](
	collection *mongo.Collection,
) IGenericMongoQueryRepository[T] {
	return &genericMongoQueryRepository[T]{
		collection: collection,
	}
}
