package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"
)

type DatabaseTestSuite struct {
	suite.Suite
	defaultDb *Source
}

type Stub struct {
	Id    primitive.ObjectID `bson:"_id"`
	Field string             `bson:"field"`
}

var (
	mongoDsn = os.Getenv("MONGO_DSN")
)

func Test_Database(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (suite *DatabaseTestSuite) SetupTest() {
	err := os.Setenv("MONGO_DIAL_TIMEOUT", "10")

	if err != nil {
		assert.FailNow(suite.T(), "Init env variable failed", "%v", err)
	}

	err = os.Setenv("MONGO_DSN", mongoDsn)

	if err != nil {
		assert.FailNow(suite.T(), "Init env variable failed", "%v", err)
	}

	db, err := NewDatabase()

	if err != nil {
		assert.FailNow(suite.T(), "New Database init failed", "%v", err)
	}

	assert.NotNil(suite.T(), db)
	assert.NotNil(suite.T(), db.connection)
	assert.IsType(suite.T(), &Options{}, db.connection)
	assert.NotNil(suite.T(), db.client)
	assert.IsType(suite.T(), &mongo.Client{}, db.client)
	assert.NotNil(suite.T(), db.collections)
	assert.Empty(suite.T(), db.collections)
	assert.NotNil(suite.T(), db.database)
	assert.IsType(suite.T(), &mongo.Database{}, db.database)
	assert.NotNil(suite.T(), db.repositoriesMu)
	assert.IsType(suite.T(), sync.Mutex{}, db.repositoriesMu)

	suite.defaultDb = db
}

func (suite *DatabaseTestSuite) TearDownTest() {
	err := suite.defaultDb.Drop()

	if err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	err = suite.defaultDb.Close()

	if err != nil {
		suite.FailNow("Database close failed", "%v", err)
	}
}

func (suite *DatabaseTestSuite) TestDatabase_StringDns_Ok() {
	conn := Options{
		Dsn: "mongodb://database_user:database_password@localhost:27017/database_name",
	}
	Host := conn.String()
	assert.NotEmpty(suite.T(), Host)
	assert.Equal(suite.T(), "mongodb://database_user:database_password@localhost:27017/database_name", Host)
}

func (suite *DatabaseTestSuite) TestDatabase_StringDnsEmpty_Ok() {
	conn := Options{
		Dsn: "some_incorrect_dns_string",
	}
	Host := conn.String()
	assert.Empty(suite.T(), Host)
}

func (suite *DatabaseTestSuite) TestDatabase_Ping_Ok() {
	err := suite.defaultDb.Ping()
	assert.NoError(suite.T(), err)
}

func (suite *DatabaseTestSuite) TestDatabase_Ping_SessionNotStart_Error() {
	sess := &Source{}
	err := sess.Ping()
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrorSessionNotInit, err)
}

func (suite *DatabaseTestSuite) TestDatabase_Collection_Ok() {
	col := suite.defaultDb.Collection("some_collection")
	assert.NotNil(suite.T(), col)
	assert.IsType(suite.T(), &mongo.Collection{}, col)
	assert.NotEmpty(suite.T(), suite.defaultDb.collections)
	assert.Len(suite.T(), suite.defaultDb.collections, 1)
	assert.Contains(suite.T(), suite.defaultDb.collections, "some_collection")
}

func (suite *DatabaseTestSuite) TestDatabase_CrudOperations_Ok() {
	var recs []*Stub

	cursor, err := suite.defaultDb.Collection("some_collection").Find(context.Background(), bson.M{})
	assert.NoError(suite.T(), err)

	err = cursor.All(context.Background(), &recs)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), recs)

	stub1 := &Stub{
		Id:    primitive.NewObjectID(),
		Field: primitive.NewObjectID().Hex(),
	}
	stub2 := &Stub{
		Id:    primitive.NewObjectID(),
		Field: primitive.NewObjectID().Hex(),
	}
	stub3 := &Stub{
		Id:    primitive.NewObjectID(),
		Field: primitive.NewObjectID().Hex(),
	}

	recsInterface := []interface{}{stub1, stub2, stub3}
	res, err := suite.defaultDb.Collection("some_collection").InsertMany(context.Background(), recsInterface)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res)

	cursor, err = suite.defaultDb.Collection("some_collection").Find(context.Background(), bson.M{})
	assert.NoError(suite.T(), err)

	err = cursor.All(context.Background(), &recs)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), recs, len(recsInterface))
	assert.Equal(suite.T(), stub1.Id, recs[0].Id)
	assert.Equal(suite.T(), stub1.Field, recs[0].Field)

	_, err = suite.defaultDb.Collection("some_collection").
		UpdateOne(
			context.Background(),
			bson.M{"_id": stub1.Id},
			bson.M{"$set": bson.M{"field": primitive.NewObjectID().Hex()}},
		)
	assert.NoError(suite.T(), err)

	var recs2 []*Stub
	cursor, err = suite.defaultDb.Collection("some_collection").Find(context.Background(), bson.M{})
	assert.NoError(suite.T(), err)

	err = cursor.All(context.Background(), &recs2)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), recs2)
	assert.Len(suite.T(), recs2, len(recsInterface))
	assert.Equal(suite.T(), recs[0].Id, recs2[0].Id)
	assert.NotEqual(suite.T(), recs[0].Field, recs2[0].Field)

	_, err = suite.defaultDb.Collection("some_collection").DeleteOne(context.Background(), bson.M{"_id": stub1.Id})
	assert.NoError(suite.T(), err)

	cursor, err = suite.defaultDb.Collection("some_collection").Find(context.Background(), bson.M{})
	assert.NoError(suite.T(), err)

	err = cursor.All(context.Background(), &recs2)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), recs2)
	assert.Len(suite.T(), recs2, len(recsInterface)-1)
}

func TestDatabase_NewDatabaseError(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	opts := []Option{
		Dsn("mongodb://database_user:database_password@incorrect_host:7777/database_name"),
		Context(ctx),
	}
	db, err := NewDatabase(opts...)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.Ping()
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func (suite *DatabaseTestSuite) TestDatabase_NewDatabaseWithOpts_Ok() {
	u, err := url.ParseRequestURI(mongoDsn)
	assert.NoError(suite.T(), err)

	err = os.Unsetenv("MONGO_DSN")
	assert.NoError(suite.T(), err)

	err = os.Unsetenv("MONGO_DIAL_TIMEOUT")
	assert.NoError(suite.T(), err)

	u.Path = "/other_db_test"

	opts := []Option{
		Dsn(u.String()),
		Mode("secondary"),
	}
	db, err := NewDatabase(opts...)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), db)
	assert.NotNil(suite.T(), db.connection)
	assert.Equal(suite.T(), u.String(), db.connection.Dsn)
	assert.IsType(suite.T(), &Options{}, db.connection)
	assert.NotNil(suite.T(), db.client)
	assert.IsType(suite.T(), &mongo.Client{}, db.client)
	assert.NotNil(suite.T(), db.collections)
	assert.Empty(suite.T(), db.collections)
	assert.NotNil(suite.T(), db.database)
	assert.IsType(suite.T(), &mongo.Database{}, db.database)
	assert.NotNil(suite.T(), db.repositoriesMu)
	assert.IsType(suite.T(), sync.Mutex{}, db.repositoriesMu)
	assert.EqualValues(suite.T(), "secondary", db.connection.Mode)
}
