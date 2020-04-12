package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SessionInterface interface {
	EndSession(context.Context)
	WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error)
	StartTransaction(...*options.TransactionOptions) error
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
	ClusterTime() bson.Raw
	AdvanceClusterTime(bson.Raw) error
	OperationTime() *primitive.Timestamp
	AdvanceOperationTime(*primitive.Timestamp) error
	Client() *mongo.Client
}

type Session struct {
	session mongo.Session
}

func (m *Session) StartTransaction(opts ...*options.TransactionOptions) error {
	return m.session.StartTransaction(opts...)
}

func (m *Session) AbortTransaction(ctx context.Context) error {
	return m.session.AbortTransaction(ctx)
}

func (m *Session) CommitTransaction(ctx context.Context) error {
	return m.session.CommitTransaction(ctx)
}

func (m *Session) WithTransaction(
	ctx context.Context,
	fn func(sessCtx mongo.SessionContext) (interface{}, error),
	opts ...*options.TransactionOptions,
) (interface{}, error) {
	return m.session.WithTransaction(ctx, fn, opts...)
}

func (m *Session) EndSession(ctx context.Context) {
	m.session.EndSession(ctx)
}

func (m *Session) ClusterTime() bson.Raw {
	return m.session.ClusterTime()
}

func (m *Session) OperationTime() *primitive.Timestamp {
	return m.session.OperationTime()
}

func (m *Session) Client() *mongo.Client {
	return m.session.Client()
}

func (m *Session) AdvanceClusterTime(in bson.Raw) error {
	return m.session.AdvanceClusterTime(in)
}

func (m *Session) AdvanceOperationTime(in *primitive.Timestamp) error {
	return m.session.AdvanceOperationTime(in)
}
