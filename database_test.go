package database

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/url"
	"os"
	"sync"
	"testing"
)

type DatabaseTestSuite struct {
	suite.Suite
	defaultDb *Source
}

type Stub struct {
	Id    bson.ObjectId `bson:"_id"`
	Field string        `bson:"field"`
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
		assert.FailNow(suite.T(), "New database init failed", "%v", err)
	}

	assert.NotNil(suite.T(), db)
	assert.NotNil(suite.T(), db.connection)
	assert.IsType(suite.T(), &Options{}, db.connection)
	assert.NotNil(suite.T(), db.session)
	assert.IsType(suite.T(), &mgo.Session{}, db.session)
	assert.NotNil(suite.T(), db.collections)
	assert.Empty(suite.T(), db.collections)
	assert.NotNil(suite.T(), db.database)
	assert.IsType(suite.T(), &mgo.Database{}, db.database)
	assert.NotNil(suite.T(), db.repositoriesMu)
	assert.IsType(suite.T(), sync.Mutex{}, db.repositoriesMu)

	suite.defaultDb = db
}

func (suite *DatabaseTestSuite) TearDownTest() {
	err := suite.defaultDb.Drop()

	if err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	suite.defaultDb.Close()
}

func (suite *DatabaseTestSuite) TestDatabase_StringDns_Ok() {
	conn := Options{
		Dsn: "mongodb://database_user:database_password@localhost:27017/database_name",
	}
	host := conn.String()
	assert.NotEmpty(suite.T(), host)
	assert.Equal(suite.T(), "mongodb://database_user:database_password@localhost:27017/database_name", host)
}

func (suite *DatabaseTestSuite) TestDatabase_StringDnsEmpty_Ok() {
	conn := Options{
		Dsn: "some_incorrect_dns_string",
	}
	host := conn.String()
	assert.Empty(suite.T(), host)
}

func (suite *DatabaseTestSuite) TestDatabase_Ping_Ok() {
	err := suite.defaultDb.Ping()
	assert.NoError(suite.T(), err)
}

func (suite *DatabaseTestSuite) TestDatabase_Ping_SessionNotStart_Error() {
	sess := &Source{}
	err := sess.Ping()
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), errorSessionNotInit, err.Error())
}

func (suite *DatabaseTestSuite) TestDatabase_Clone_Ok() {
	sess := suite.defaultDb.Clone()
	assert.NotNil(suite.T(), sess)
	assert.NotEqual(suite.T(), suite.defaultDb, sess)
}

func (suite *DatabaseTestSuite) TestDatabase_Collection_Ok() {
	col := suite.defaultDb.Collection("some_collection")
	assert.NotNil(suite.T(), col)
	assert.IsType(suite.T(), &mgo.Collection{}, col)
	assert.NotEmpty(suite.T(), suite.defaultDb.collections)
	assert.Len(suite.T(), suite.defaultDb.collections, 1)
	assert.Contains(suite.T(), suite.defaultDb.collections, "some_collection")
}

func (suite *DatabaseTestSuite) TestDatabase_CrudOperations_Ok() {
	var recs []*Stub

	err := suite.defaultDb.Collection("some_collection").Find(bson.M{}).All(&recs)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), recs)

	stub1 := &Stub{
		Id:    bson.NewObjectId(),
		Field: bson.NewObjectId().Hex(),
	}
	stub2 := &Stub{
		Id:    bson.NewObjectId(),
		Field: bson.NewObjectId().Hex(),
	}
	stub3 := &Stub{
		Id:    bson.NewObjectId(),
		Field: bson.NewObjectId().Hex(),
	}

	recsInterface := []interface{}{stub1, stub2, stub3}
	err = suite.defaultDb.Collection("some_collection").Insert(recsInterface...)
	assert.NoError(suite.T(), err)

	err = suite.defaultDb.Collection("some_collection").Find(bson.M{}).All(&recs)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), recs)
	assert.Len(suite.T(), recs, len(recsInterface))
	assert.Equal(suite.T(), stub1.Id, recs[0].Id)
	assert.Equal(suite.T(), stub1.Field, recs[0].Field)

	stub1.Field = bson.NewObjectId().Hex()
	err = suite.defaultDb.Collection("some_collection").UpdateId(stub1.Id, stub1)
	assert.NoError(suite.T(), err)

	var recs2 []*Stub
	err = suite.defaultDb.Collection("some_collection").Find(bson.M{}).All(&recs2)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), recs2)
	assert.Len(suite.T(), recs2, len(recsInterface))
	assert.Equal(suite.T(), recs[0].Id, recs2[0].Id)
	assert.NotEqual(suite.T(), recs[0].Field, recs2[0].Field)

	err = suite.defaultDb.Collection("some_collection").RemoveId(stub1.Id)
	assert.NoError(suite.T(), err)

	err = suite.defaultDb.Collection("some_collection").Find(bson.M{}).All(&recs2)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), recs2)
	assert.Len(suite.T(), recs2, len(recsInterface)-1)
}

func (suite *DatabaseTestSuite) TestDatabase_EnvVariablesParse_Error() {
	err := os.Setenv("MONGO_DIAL_TIMEOUT", "qwerty")
	assert.NoError(suite.T(), err)

	db, err := NewDatabase()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), db)
	assert.Regexp(suite.T(), "MONGO_DIAL_TIMEOUT", err.Error())
}

func (suite *DatabaseTestSuite) TestDatabase_NewDatabaseError() {
	err := os.Setenv("MONGO_DSN", "mongodb://database_user:database_password@incorrect_host:7777/database_name")
	assert.NoError(suite.T(), err)

	err = os.Setenv("MONGO_DIAL_TIMEOUT", "1")
	assert.NoError(suite.T(), err)

	db, err := NewDatabase()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), db)
	assert.Equal(suite.T(), "no reachable servers", err.Error())

	err = os.Unsetenv("MONGO_DSN")
	assert.NoError(suite.T(), err)
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
		DialTimeout(1),
		Mode(mgo.Secondary),
	}
	db, err := NewDatabase(opts...)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), db)
	assert.NotNil(suite.T(), db.connection)
	assert.Equal(suite.T(), u.String(), db.connection.Dsn)
	assert.Equal(suite.T(), int64(1), db.connection.DialTimeout)
	assert.IsType(suite.T(), &Options{}, db.connection)
	assert.NotNil(suite.T(), db.session)
	assert.IsType(suite.T(), &mgo.Session{}, db.session)
	assert.NotNil(suite.T(), db.collections)
	assert.Empty(suite.T(), db.collections)
	assert.NotNil(suite.T(), db.database)
	assert.IsType(suite.T(), &mgo.Database{}, db.database)
	assert.NotNil(suite.T(), db.repositoriesMu)
	assert.IsType(suite.T(), sync.Mutex{}, db.repositoriesMu)
	assert.Equal(suite.T(), mgo.Secondary, db.connection.Mode)
}
