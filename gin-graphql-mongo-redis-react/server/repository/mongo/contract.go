package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IGenericMongoCommandRepository[T any] interface {
	InsertData(
		ctx context.Context,
		params bson.M,
	) error
	InsertBatchData(
		ctx context.Context,
		body []interface{},
	) error
	UpdateData(
		ctx context.Context,
		filters, params bson.M,
	) error
	UpdateBatchData(
		ctx context.Context,
		filters, params bson.M,
	) error
	CreateOrUpdateData(
		ctx context.Context,
		filters, params bson.M,
	) error
	DeleteData(
		ctx context.Context,
		filters bson.M,
	) error
	DeleteBatchData(
		ctx context.Context,
		filters bson.M,
	) error
}

type IGenericMongoQueryRepository[T any] interface {
	CountData(
		ctx context.Context,
		pipeline mongo.Pipeline,
	) int64
	AggregateData(
		ctx context.Context,
		pipeline mongo.Pipeline,
	) ([]*T, error)
	FindData(
		ctx context.Context,
		filter bson.M,
	) (*T, error)
	DistinctData(
		ctx context.Context,
		filter bson.M,
		field string,
	) ([]interface{}, error)
}
