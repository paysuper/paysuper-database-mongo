package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionInterface interface {
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (CursorInterface, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultInterface
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) SingleResultInterface
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) SingleResultInterface
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) SingleResultInterface
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	Indexes() mongo.IndexView
}

type SingleResultInterface interface {
	Decode(v interface{}) error
	DecodeBytes() (bson.Raw, error)
	Err() error
}

type Collection struct {
	collection *mongo.Collection
}

type SingleResult struct {
	singleResult *mongo.SingleResult
}

func (m *Collection) Aggregate(
	ctx context.Context,
	pipeline interface{},
	opts ...*options.AggregateOptions,
) (CursorInterface, error) {
	cursor, err := m.collection.Aggregate(ctx, pipeline, opts...)

	if err != nil {
		return nil, err
	}

	return &Cursor{cursor: cursor}, nil
}

func (m *Collection) CountDocuments(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (int64, error) {
	return m.collection.CountDocuments(ctx, filter, opts...)
}

func (m *Collection) DeleteMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {
	return m.collection.DeleteMany(ctx, filter, opts...)
}

func (m *Collection) DeleteOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {
	return m.collection.DeleteOne(ctx, filter, opts...)
}

func (m *Collection) Distinct(
	ctx context.Context,
	fieldName string,
	filter interface{},
	opts ...*options.DistinctOptions,
) ([]interface{}, error) {
	return m.collection.Distinct(ctx, fieldName, filter, opts...)
}

func (m *Collection) Find(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (CursorInterface, error) {
	cursor, err := m.collection.Find(ctx, filter, opts...)

	if err != nil {
		return nil, err
	}

	return &Cursor{cursor: cursor}, nil
}

func (m *Collection) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) SingleResultInterface {
	result := m.collection.FindOne(ctx, filter, opts...)
	return &SingleResult{singleResult: result}
}

func (m *Collection) FindOneAndDelete(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneAndDeleteOptions,
) SingleResultInterface {
	result := m.collection.FindOneAndDelete(ctx, filter, opts...)
	return &SingleResult{singleResult: result}
}

func (m *Collection) FindOneAndReplace(
	ctx context.Context,
	filter interface{},
	replacement interface{},
	opts ...*options.FindOneAndReplaceOptions,
) SingleResultInterface {
	result := m.collection.FindOneAndReplace(ctx, filter, replacement, opts...)
	return &SingleResult{singleResult: result}
}

func (m *Collection) FindOneAndUpdate(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.FindOneAndUpdateOptions,
) SingleResultInterface {
	result := m.collection.FindOneAndUpdate(ctx, filter, update, opts...)
	return &SingleResult{singleResult: result}
}

func (m *Collection) InsertMany(
	ctx context.Context,
	documents []interface{},
	opts ...*options.InsertManyOptions,
) (*mongo.InsertManyResult, error) {
	return m.collection.InsertMany(ctx, documents, opts...)
}

func (m *Collection) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...*options.InsertOneOptions,
) (*mongo.InsertOneResult, error) {
	return m.collection.InsertOne(ctx, document, opts...)
}

func (m *Collection) ReplaceOne(
	ctx context.Context,
	filter interface{},
	replacement interface{},
	opts ...*options.ReplaceOptions,
) (*mongo.UpdateResult, error) {
	return m.collection.ReplaceOne(ctx, filter, replacement, opts...)
}

func (m *Collection) UpdateMany(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	return m.collection.UpdateMany(ctx, filter, update, opts...)
}

func (m *Collection) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	return m.collection.UpdateOne(ctx, filter, update, opts...)
}

func (m *Collection) BulkWrite(
	ctx context.Context,
	models []mongo.WriteModel,
	opts ...*options.BulkWriteOptions,
) (*mongo.BulkWriteResult, error) {
	return m.collection.BulkWrite(ctx, models, opts...)
}

func (m *Collection) Indexes() mongo.IndexView {
	return m.collection.Indexes()
}

func (m *SingleResult) Decode(v interface{}) error {
	return m.singleResult.Decode(v)
}

func (m *SingleResult) DecodeBytes() (bson.Raw, error) {
	return m.singleResult.DecodeBytes()
}

func (m *SingleResult) Err() error {
	return m.singleResult.Err()
}
