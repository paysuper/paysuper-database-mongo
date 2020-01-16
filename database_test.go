package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"
)

type DatabaseTestSuite struct {
	suite.Suite
	defaultDb SourceInterface
}

var (
	mongoDsn = os.Getenv("MONGO_DSN")
)

func Test_Database(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (suite *DatabaseTestSuite) SetupTest() {
	err := os.Setenv("MONGO_DSN", mongoDsn)
	if err != nil {
		assert.FailNow(suite.T(), "Init env variable failed", "%v", err)
	}

	db, err := NewDatabase()
	if err != nil {
		assert.FailNow(suite.T(), "New Database init failed", "%v", err)
	}

	assert.NotNil(suite.T(), db)
	assert.Implements(suite.T(), (*SourceInterface)(nil), db)
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
	err := suite.defaultDb.Ping(context.TODO())
	assert.NoError(suite.T(), err)
}

func (suite *DatabaseTestSuite) TestDatabase_Ping_SessionNotStart_Error() {
	sess := &Source{}
	err := sess.Ping(context.TODO())
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrorSessionNotInit, err)
}

func (suite *DatabaseTestSuite) TestDatabase_Collection_Ok() {
	col := suite.defaultDb.Collection("some_collection")
	assert.NotNil(suite.T(), col)
	assert.Implements(suite.T(), (*CollectionInterface)(nil), col)

	db, ok := suite.defaultDb.(*Source)
	assert.True(suite.T(), ok)
	assert.NotEmpty(suite.T(), db.collections)
	assert.Len(suite.T(), db.collections, 1)
	assert.Contains(suite.T(), db.collections, "some_collection")
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

	err = db.Ping(nil)
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
	db0, err := NewDatabase(opts...)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), db0)

	db, ok := db0.(*Source)
	assert.True(suite.T(), ok)
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

func (suite *DatabaseTestSuite) TestDatabase_NewDatabase_Error() {
	opts := []Option{
		Dsn("some_incorrect_dns_string"),
	}
	db, err := NewDatabase(opts...)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrorProtocolNotFound, err)
	assert.Nil(suite.T(), db)
}

func (suite *DatabaseTestSuite) TestDatabase_NewDatabase_DsnDatabaseNotFound_Error() {
	opts := []Option{
		Dsn("mongodb://localhost:27017/"),
	}
	db, err := NewDatabase(opts...)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrorDatabaseNotFound, err)
	assert.Nil(suite.T(), db)
}

func (suite *DatabaseTestSuite) TestDatabase_NewDatabase_IncorrectReadPref_Error() {
	opts := []Option{
		Dsn(mongoDsn),
		Mode("some_incorrect_value"),
	}
	db, err := NewDatabase(opts...)
	assert.Error(suite.T(), err)
	assert.Regexp(suite.T(), "unknown read preference", err.Error())
	assert.Nil(suite.T(), db)
}
