package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type genericMongoCommandRepository[T any] struct {
	collection *mongo.Collection
}

func (repository genericMongoCommandRepository[T]) InsertData(
	ctx context.Context,
	params bson.M,
) error {
	if _, err := repository.collection.InsertOne(
		ctx, params,
	); err != nil {
		return err
	}
	return nil
}

func (repository genericMongoCommandRepository[T]) InsertBatchData(
	ctx context.Context,
	body []interface{},
) error {
	if _, err := repository.collection.InsertMany(ctx, body); err != nil {
		return err
	}
	return nil
}

func (repository genericMongoCommandRepository[T]) UpdateData(
	ctx context.Context,
	filters, params bson.M,
) error {
	return repository.collection.FindOneAndUpdate(
		ctx, filters, params,
	).Err()
}

func (repository genericMongoCommandRepository[T]) UpdateBatchData(
	ctx context.Context,
	filters, params bson.M,
) error {
	if _, err := repository.collection.UpdateMany(
		ctx, filters, params,
	); err != nil {
		return err
	}
	return nil
}

func (repository genericMongoCommandRepository[T]) CreateOrUpdateData(
	ctx context.Context,
	filters, params bson.M,
) error {
	opts := options.Update().SetUpsert(true)
	_, err := repository.collection.UpdateOne(ctx, filters, params, opts)
	if err != nil {
		return err
	}
	return nil
}

func (repository genericMongoCommandRepository[T]) DeleteData(
	ctx context.Context,
	filters bson.M,
) error {
	if _, err := repository.collection.DeleteOne(
		ctx, filters, nil,
	); err != nil {
		return err
	}
	return nil
}

func (repository genericMongoCommandRepository[T]) DeleteBatchData(
	ctx context.Context,
	filters bson.M,
) error {
	if _, err := repository.collection.DeleteMany(ctx, filters); err != nil {
		return err
	}
	return nil
}

func NewGenericMongoCommandRepository[T any](
	collection *mongo.Collection,
) IGenericMongoCommandRepository[T] {
	return &genericMongoCommandRepository[T]{
		collection: collection,
	}
}
