package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type CursorInterface interface {
	All(ctx context.Context, results interface{}) error
	Close(ctx context.Context) error
	Decode(val interface{}) error
	Err() error
	ID() int64
	Next(ctx context.Context) bool
	TryNext(ctx context.Context) bool
}

type Cursor struct {
	cursor *mongo.Cursor
}

func (m *Cursor) All(ctx context.Context, results interface{}) error {
	return m.cursor.All(ctx, results)
}

func (m *Cursor) Close(ctx context.Context) error {
	return m.cursor.Close(ctx)
}

func (m *Cursor) Decode(val interface{}) error {
	return m.cursor.Decode(val)
}

func (m *Cursor) Err() error {
	return m.cursor.Err()
}

func (m *Cursor) ID() int64 {
	return m.cursor.ID()
}

func (m *Cursor) Next(ctx context.Context) bool {
	return m.cursor.Next(ctx)
}

func (m *Cursor) TryNext(ctx context.Context) bool {
	return m.cursor.Next(ctx)
}
