package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
)

var (
	stubCollection = "stubs"
)

type Stub struct {
	FieldString string  `bson:"field_string"`
	FieldFloat  float64 `bson:"field_float"`
}

type CollectionTestSuite struct {
	suite.Suite
	defaultDb SourceInterface
	stubs     []interface{}
}

func Test_Collection(t *testing.T) {
	suite.Run(t, new(CollectionTestSuite))
}

func (suite *CollectionTestSuite) SetupTest() {
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

func (suite *CollectionTestSuite) TearDownTest() {
	err := suite.defaultDb.Drop()
	if err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	err = suite.defaultDb.Close()
	if err != nil {
		suite.FailNow("Database close failed", "%v", err)
	}
}

func (suite *CollectionTestSuite) TestCollection_Aggregate_Ok() {
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":    "$field_string",
				"amount": bson.M{"$sum": "$field_float"},
			},
		},
	}
	cursor, err := suite.defaultDb.Collection(stubCollection).Aggregate(context.TODO(), pipeline)
	assert.NoError(suite.T(), err)

	var result []struct {
		Id     string  `bson:"_id"`
		Amount float64 `bson:"amount"`
	}
	err = cursor.All(context.TODO(), &result)
	assert.Len(suite.T(), result, 4)

	for _, v := range result {
		if v.Id == "value1" {
			assert.EqualValues(suite.T(), v.Amount, 60)
		}

		if v.Id == "value2" {
			assert.EqualValues(suite.T(), v.Amount, 20)
		}

		if v.Id == "value3" {
			assert.EqualValues(suite.T(), v.Amount, 150)
		}

		if v.Id == "value4" {
			assert.EqualValues(suite.T(), v.Amount, 10)
		}
	}

	err = cursor.Close(context.TODO())
	assert.NoError(suite.T(), err)
}

func (suite *CollectionTestSuite) TestCollection_Aggregate_Error() {
	pipeline := []bson.M{
		{
			"$unknownFn": bson.M{
				"_id":    "$field_string",
				"amount": bson.M{"$sum": "$field_float"},
			},
		},
	}
	cursor, err := suite.defaultDb.Collection(stubCollection).Aggregate(context.TODO(), pipeline)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), cursor)
	tErr, ok := err.(mongo.CommandError)
	assert.True(suite.T(), ok)
	assert.EqualValues(suite.T(), 40324, tErr.Code)
	assert.Regexp(suite.T(), "\\$unknownFn", tErr.Message)
}

func (suite *CollectionTestSuite) TestCollection_CountDocuments_Ok() {
	count, err := suite.defaultDb.Collection(stubCollection).CountDocuments(context.TODO(), bson.M{})
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), count, len(suite.stubs))
}

func (suite *CollectionTestSuite) TestCollection_DeleteMany_Ok() {
	res, err := suite.defaultDb.Collection(stubCollection).DeleteMany(context.TODO(), bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), res.DeletedCount, 3)
}

func (suite *CollectionTestSuite) TestCollection_DeleteOne_Ok() {
	ctx := context.TODO()

	res, err := suite.defaultDb.Collection(stubCollection).DeleteOne(ctx, bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), res.DeletedCount, 1)

	cursor, err := suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)

	var result []*Stub
	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
}

func (suite *CollectionTestSuite) TestCollection_Distinct_Ok() {
	res, err := suite.defaultDb.Collection(stubCollection).Distinct(context.TODO(), "field_string", bson.M{})
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), res, 4)
}

func (suite *CollectionTestSuite) TestCollection_Find_Error() {
	filter := bson.M{"field_string": bson.M{"$unknownFn": "val"}}
	cursor, err := suite.defaultDb.Collection(stubCollection).Find(context.TODO(), filter)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), cursor)
	tErr, ok := err.(mongo.CommandError)
	assert.True(suite.T(), ok)
	assert.EqualValues(suite.T(), 2, tErr.Code)
	assert.Regexp(suite.T(), "\\$unknownFn", tErr.Message)
}

func (suite *CollectionTestSuite) TestCollection_FindOne_Ok() {
	var res *Stub
	err := suite.defaultDb.Collection(stubCollection).FindOne(context.TODO(), bson.M{"field_string": "value1"}).Decode(&res)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *CollectionTestSuite) TestCollection_FindOneAndDelete_Ok() {
	ctx := context.TODO()

	err := suite.defaultDb.Collection(stubCollection).FindOneAndDelete(ctx, bson.M{"field_string": "value1"}).Err()
	assert.NoError(suite.T(), err)

	cursor, err := suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)

	var result []*Stub
	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
}

func (suite *CollectionTestSuite) TestCollection_FindOneAndReplace_Ok() {
	newVal := &Stub{FieldString: "value5", FieldFloat: 111}
	ctx := context.TODO()

	err := suite.defaultDb.Collection(stubCollection).FindOneAndReplace(ctx, bson.M{"field_string": "value1"}, newVal).Err()
	assert.NoError(suite.T(), err)

	cursor, err := suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)

	var result []*Stub
	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	cursor, err = suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value5"})
	assert.NoError(suite.T(), err)

	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
}

func (suite *CollectionTestSuite) TestCollection_FindOneAndUpdate_Ok() {
	ctx := context.TODO()

	err := suite.defaultDb.Collection(stubCollection).
		FindOneAndUpdate(
			ctx,
			bson.M{"field_string": "value1"},
			bson.M{"$set": bson.M{"field_string": "value6"}},
		).Err()
	assert.NoError(suite.T(), err)

	cursor, err := suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)

	var result []*Stub
	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	cursor, err = suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value6"})
	assert.NoError(suite.T(), err)

	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
}

func (suite *CollectionTestSuite) TestCollection_InsertOne_Ok() {
	ctx := context.TODO()

	err := suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value6"}).Err()
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), mongo.ErrNoDocuments, err)

	doc := &Stub{FieldString: "value6", FieldFloat: 111}
	res, err := suite.defaultDb.Collection(stubCollection).InsertOne(ctx, doc)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res.InsertedID)

	var res1 *Stub
	err = suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value6"}).Decode(&res1)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res1)
	assert.Equal(suite.T(), doc, res1)
}

func (suite *CollectionTestSuite) TestCollection_ReplaceOne_Ok() {
	ctx := context.TODO()
	var res *Stub

	err := suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value4"}).Decode(&res)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), res.FieldString, "value4")

	err = suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value6"}).Err()
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), mongo.ErrNoDocuments, err)

	doc := &Stub{FieldString: "value6", FieldFloat: 111}
	res1, err := suite.defaultDb.Collection(stubCollection).ReplaceOne(ctx, bson.M{"field_string": "value4"}, doc)
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), res1.MatchedCount, 1)
	assert.EqualValues(suite.T(), res1.ModifiedCount, 1)

	err = suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value6"}).Decode(&res)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), doc, res)

	err = suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value4"}).Err()
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), err, mongo.ErrNoDocuments)
}

func (suite *CollectionTestSuite) TestCollection_UpdateMany_Ok() {
	ctx := context.TODO()

	cursor, err := suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)

	var result []*Stub
	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)

	cursor, err = suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value6"})
	assert.NoError(suite.T(), err)
	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)

	res, err := suite.defaultDb.Collection(stubCollection).UpdateMany(
		ctx,
		bson.M{"field_string": "value1"},
		bson.M{"$set": bson.M{"field_string": "value6"}},
	)
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), res.MatchedCount, 3)
	assert.EqualValues(suite.T(), res.ModifiedCount, 3)

	cursor, err = suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value6"})
	assert.NoError(suite.T(), err)
	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)

	cursor, err = suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)
	err = cursor.All(ctx, &result)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

func (suite *CollectionTestSuite) TestCollection_UpdateOne_Ok() {
	ctx := context.TODO()
	var stub *Stub

	err := suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value4"}).Decode(&stub)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stub)

	err = suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value6"}).Err()
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), err, mongo.ErrNoDocuments)

	res, err := suite.defaultDb.Collection(stubCollection).UpdateOne(
		ctx,
		bson.M{"field_string": "value4"},
		bson.M{"$set": bson.M{"field_string": "value6"}},
	)
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), res.MatchedCount, 1)
	assert.EqualValues(suite.T(), res.ModifiedCount, 1)

	err = suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value6"}).Decode(&stub)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stub)

	err = suite.defaultDb.Collection(stubCollection).FindOne(ctx, bson.M{"field_string": "value4"}).Err()
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), err, mongo.ErrNoDocuments)
}

func (suite *CollectionTestSuite) TestCollection_BulkWrite_Ok() {
	ctx := context.TODO()
	models := []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"field_string": "value1"}).
			SetUpdate(bson.M{"$set": bson.M{"field_float": 1}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"field_string": "value2"}).
			SetUpdate(bson.M{"$set": bson.M{"field_float": 2}}),
	}

	var results []*Stub

	cursor, err := suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)
	err = cursor.All(ctx, &results)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), results)

	for _, v := range results {
		assert.NotEqual(suite.T(), v.FieldFloat, 1)
	}

	cursor, err = suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value2"})
	assert.NoError(suite.T(), err)
	err = cursor.All(ctx, &results)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), results)

	for _, v := range results {
		assert.NotEqual(suite.T(), v.FieldFloat, 2)
	}

	res, err := suite.defaultDb.Collection(stubCollection).BulkWrite(ctx, models)
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), res.MatchedCount, 5)
	assert.EqualValues(suite.T(), res.ModifiedCount, 5)

	cursor, err = suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value1"})
	assert.NoError(suite.T(), err)
	err = cursor.All(ctx, &results)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), results)

	for _, v := range results {
		assert.EqualValues(suite.T(), v.FieldFloat, 1)
	}

	cursor, err = suite.defaultDb.Collection(stubCollection).Find(ctx, bson.M{"field_string": "value2"})
	assert.NoError(suite.T(), err)
	err = cursor.All(ctx, &results)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), results)

	for _, v := range results {
		assert.EqualValues(suite.T(), v.FieldFloat, 2)
	}
}

func (suite *CollectionTestSuite) TestCollection_Indexes_Ok() {
	res := suite.defaultDb.Collection(stubCollection).Indexes()
	assert.NotNil(suite.T(), res)
}

func (suite *CollectionTestSuite) TestCollection_SingleResult_DecodeBytes_Ok() {
	res := suite.defaultDb.Collection(stubCollection).FindOne(context.TODO(), bson.M{"field_string": "value4"})
	assert.NoError(suite.T(), res.Err())

	var firstDecodedResult bson.Raw
	err := res.Decode(&firstDecodedResult)
	assert.NoError(suite.T(), err)
	secondDecodedResult, err := res.DecodeBytes()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), firstDecodedResult, secondDecodedResult)
}
