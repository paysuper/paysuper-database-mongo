package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"testing"
)

type CursorTestSuite struct {
	suite.Suite
	defaultDb *Source
	stubs     []interface{}
}

func Test_Cursor(t *testing.T) {
	suite.Run(t, new(CursorTestSuite))
}

func (suite *CursorTestSuite) SetupTest() {
	err := os.Setenv("MONGO_DSN", mongoDsn)
	if err != nil {
		assert.FailNow(suite.T(), "Init env variable failed", "%v", err)
	}

	suite.defaultDb, err = NewDatabase()
	if err != nil {
		assert.FailNow(suite.T(), "New Database init failed", "%v", err)
	}

	assert.NotNil(suite.T(), suite.defaultDb)

	suite.stubs = []interface{}{
		&Stub{FieldString: "value1", FieldFloat: 10},
		&Stub{FieldString: "value1", FieldFloat: 20},
		&Stub{FieldString: "value1", FieldFloat: 30},
		&Stub{FieldString: "value2", FieldFloat: 10},
		&Stub{FieldString: "value2", FieldFloat: 10},
		&Stub{FieldString: "value3", FieldFloat: 100},
		&Stub{FieldString: "value3", FieldFloat: 20},
		&Stub{FieldString: "value3", FieldFloat: 30},
		&Stub{FieldString: "value4", FieldFloat: 10},
	}
	res, err := suite.defaultDb.Collection(stubCollection).InsertMany(context.TODO(), suite.stubs)
	if err != nil {
		assert.FailNow(suite.T(), "Add stub data to collection failed", "%v", err)
	}

	assert.Len(suite.T(), res.InsertedIDs, len(suite.stubs))
}

func (suite *CursorTestSuite) TearDownTest() {
	err := suite.defaultDb.Drop()
	if err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	err = suite.defaultDb.Close()
	if err != nil {
		suite.FailNow("Database close failed", "%v", err)
	}
}

func (suite *CursorTestSuite) TestCursor_Ok() {
	ctx := context.TODO()
	pipeline := []bson.M{
		{
			"$match": bson.M{"field_string": "value1"},
		},
		{
			"$group": bson.M{
				"_id":    "$field_string",
				"amount": bson.M{"$sum": "$field_float"},
			},
		},
	}
	cursor, err := suite.defaultDb.Collection(stubCollection).Aggregate(ctx, pipeline)
	assert.NoError(suite.T(), err)

	var result struct {
		Id     string  `bson:"_id"`
		Amount float64 `bson:"amount"`
	}

	assert.NoError(suite.T(), cursor.Err())
	assert.True(suite.T(), cursor.Next(ctx))
	assert.EqualValues(suite.T(), cursor.ID(), 0)
	err = cursor.Decode(&result)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), result.Id, "value1")
	assert.EqualValues(suite.T(), result.Amount, 60)
	assert.False(suite.T(), cursor.TryNext(ctx))

	err = cursor.Close(ctx)
	assert.NoError(suite.T(), err)
}
