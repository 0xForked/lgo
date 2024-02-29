package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	mongorepo "learn.go/gin-graphql-mongo-redis/server/repository/mongo"
)

type mockStruct struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func Test_AggregateData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	id := primitive.NewObjectID().String()
	mt.Run("test aggregate should return data", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1,
			"job.data",
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: id},
			})
		killCursors := mtest.CreateCursorResponse(
			0, "job.data", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data, err := repo.AggregateData(context.TODO(), mongo.Pipeline{})
		assert.Nil(mt, err)
		assert.NotNil(mt, data)
		assert.Equal(mt, data[0].ID, id)
	})
}
func Test_AggregateData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	id := primitive.NewObjectID().String()
	mt.Run("test aggregate should return data and error on decode", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1,
			"job.data",
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: id},
			})
		second := mtest.CreateCursorResponse(1,
			"job.data",
			mtest.NextBatch,
			bson.D{
				{Key: "_id", Value: 123},
			})
		mt.AddMockResponses(first, second)
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data, err := repo.AggregateData(context.TODO(), mongo.Pipeline{})
		assert.NotNil(mt, err)
		assert.Nil(mt, data)
		assert.Equal(mt, err.Error(), "error decoding key _id: cannot decode 32-bit integer into a string type")
	})
	mt.Run("test should error aggregate", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    123,
			Message: "READ_FAILED",
		}))
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data, err := repo.AggregateData(context.TODO(), mongo.Pipeline{})
		assert.Nil(mt, data)
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), "write command error: [{write errors: [{READ_FAILED}]}, {<nil>}]")
	})
}

func Test_CountData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test count should return data from aggregate", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1,
				"count.data",
				mtest.FirstBatch,
				bson.D{
					{Key: "count", Value: 12},
				}),
		)
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data := repo.CountData(context.TODO(), nil)
		assert.NotNil(mt, data)
		assert.NotZero(mt, data)
	})
}
func Test_CountData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test count doc should error", func(mt *mtest.T) {
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data := repo.CountData(context.TODO(), nil)
		assert.Zero(mt, data)
	})

	mt.Run("test count should error from aggregate", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    123,
			Message: "READ_FAILED",
		}))
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data := repo.CountData(context.TODO(), nil)
		assert.Zero(mt, data)
	})

	mt.Run("test count should error from decode", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1,
			"count.data",
			mtest.FirstBatch,
			bson.D{
				{Key: "count", Value: "12"},
			}))
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data := repo.CountData(context.TODO(), nil)
		assert.Zero(mt, data)
	})

	mt.Run("test count should error from cursor", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(50, "foo.bar", mtest.FirstBatch),
			mtest.CreateSuccessResponse(),
		)
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data := repo.CountData(context.TODO(), nil)
		assert.Zero(mt, data)
	})
}

func Test_FindData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test find should return data", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"job.data",
				mtest.FirstBatch,
				bson.D{
					{"_id", "123"},
				},
			),
		)
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data, err := repo.FindData(context.TODO(), bson.M{"_id": "123"})
		assert.Nil(mt, err)
		assert.NotNil(mt, data)
		assert.Equal(mt, data.ID, "123")
	})
}
func Test_FindData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test find return error on query", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    123,
			Message: "READ_FAILED",
		}))
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data, err := repo.FindData(context.TODO(), bson.M{"_id": "123"})
		assert.Nil(mt, data)
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), "write command error: [{write errors: [{READ_FAILED}]}, {<nil>}]")
	})
}

func Test_DistinctData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test distinct return items", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"values", bson.A{"team1"}},
		})
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data, err := repo.DistinctData(context.TODO(), bson.M{}, "")
		assert.Nil(t, err)
		assert.NotNil(t, data)
	})
}
func Test_DistinctData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test distinct return error on query", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    123,
			Message: "READ_FAILED",
		}))
		repo := mongorepo.NewGenericMongoQueryRepository[mockStruct](mt.Coll)
		data, err := repo.DistinctData(context.TODO(), bson.M{}, "")
		assert.NotNil(t, err)
		assert.Nil(t, data)
		assert.Equal(mt, err.Error(), "write command error: [{write errors: [{READ_FAILED}]}, {<nil>}]")
	})
}

func Test_InsertData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test insert return nil", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(
				bson.E{
					Key:   "_id",
					Value: "1234-abcd-5678",
				},
			),
		)
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.InsertData(context.TODO(), bson.M{"_id": "1234-abcd-5678"})
		assert.Nil(mt, err)
	})
}
func Test_InsertData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test insert return error", func(mt *mtest.T) {
		errRes := mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    123,
			Message: "FAILED_WRITE",
		})
		mt.AddMockResponses(errRes)
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.InsertData(context.TODO(), bson.M{"_id": "1234-abcd-5678"})
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), "write exception: write errors: [FAILED_WRITE]")
	})
}

func Test_InsertBatchData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test insert return nil", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(
				bson.E{
					Key:   "_id",
					Value: "1234-abcd-5678",
				},
			),
		)
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.InsertBatchData(context.TODO(), []interface{}{
			bson.M{"_id": "1234-abcd-5678"},
			bson.M{"_id": "1234-abcd-9876"},
		})
		assert.Nil(mt, err)
	})
}
func Test_InsertBatchData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test insert return error", func(mt *mtest.T) {
		errRes := mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    123,
			Message: "FAILED_WRITE",
		})
		mt.AddMockResponses(errRes)
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.InsertBatchData(context.TODO(), []interface{}{
			bson.M{"_id": "1234-abcd-5678"},
			bson.M{"_id": "1234-abcd-9876"},
		})
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), "bulk write exception: write errors: [FAILED_WRITE]")
	})
}

func Test_UpdateData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test update should success", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"value", bson.D{
				{"_id", "1234-abcd-5678"},
				{"name", "lorem"},
				{"updated_at", time.Now()},
			}},
		})
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.UpdateData(context.TODO(), bson.M{
			"_id": "1234-abcd-5678",
		}, bson.M{
			"$set": bson.M{
				"name":       "lorem",
				"updated_at": time.Now(),
			},
		})
		assert.Nil(mt, err)
	})
}
func Test_UpdateData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test update one should failed param", func(mt *mtest.T) {
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.UpdateData(context.TODO(), bson.M{}, bson.M{})
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), "update document must have at least one element")
	})
	mt.Run("test update one should failed find", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{
				Key:   "_id",
				Value: "1234-abcd-5678",
			},
		))
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.UpdateData(context.TODO(), bson.M{
			"_id": "1234-abcd-5678",
		}, bson.M{
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		})
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), mongo.ErrNoDocuments.Error())
	})
}

func Test_UpdateBatchData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test insert return nil", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(
				bson.E{
					Key:   "_id",
					Value: "1234-abcd-5678",
				},
			),
		)
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.UpdateBatchData(context.TODO(), bson.M{
			"_id": "1234-abcd-5678",
		}, bson.M{
			"$set": bson.M{
				"name":       "lorem",
				"updated_at": time.Now(),
			},
		})
		assert.Nil(mt, err)
	})
}
func Test_UpdateBatchData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test insert return error", func(mt *mtest.T) {
		errRes := mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    123,
			Message: "FAILED_WRITE",
		})
		mt.AddMockResponses(errRes)
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.UpdateBatchData(context.TODO(), bson.M{
			"_id": "1234-abcd-5678",
		}, bson.M{
			"$set": bson.M{
				"name":       "lorem",
				"updated_at": time.Now(),
			},
		})
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), "write exception: write errors: [FAILED_WRITE]")
	})
}

func Test_CreateOrUpdateData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test update should success", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"value", bson.D{
				{"_id", "1234-abcd-5678"},
				{"name", "lorem"},
				{"updated_at", time.Now()},
			}},
		})
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.CreateOrUpdateData(context.TODO(), bson.M{
			"_id": "1234-abcd-5678",
		}, bson.M{
			"$set": bson.M{
				"_id":        "1234-abcd-5678",
				"name":       "lorem",
				"updated_at": time.Now(),
			},
		})
		assert.Nil(mt, err)
	})
}
func Test_CreateOrUpdateData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test update one should failed find", func(mt *mtest.T) {
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.CreateOrUpdateData(context.TODO(), bson.M{
			"_id": "1234-abcd-5678",
		}, bson.M{
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		})
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), "no responses remaining")
	})
}

func Test_DeleteData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test hard delete", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"acknowledged", true},
			{"n", 1},
		})
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.DeleteData(context.TODO(), bson.M{"_id": "1234-abcd-5678"})
		assert.Nil(mt, err)
	})
}
func Test_DeleteData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test hard delete", func(mt *mtest.T) {
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.DeleteData(context.TODO(), bson.M{"_id": "1234-abcd-5678"})
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), "no responses remaining")
	})
}

func Test_DeleteBatchData_ShouldSuccess(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test hard delete", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"acknowledged", true},
			{"n", 1},
		})
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.DeleteBatchData(context.TODO(), bson.M{"tenant_id": "1234-abcd-5678"})
		assert.Nil(mt, err)
	})
}
func Test_DeleteBatchData_ShouldError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test hard delete", func(mt *mtest.T) {
		repo := mongorepo.NewGenericMongoCommandRepository[mockStruct](mt.Coll)
		err := repo.DeleteBatchData(context.TODO(), bson.M{"tenant_id": "1234-abcd-5678"})
		assert.NotNil(mt, err)
		assert.Equal(mt, err.Error(), "no responses remaining")
	})
}
