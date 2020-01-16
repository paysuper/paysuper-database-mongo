package main

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	database "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"gopkg.in/paysuper/paysuper-database-mongo.v2/mocks"
	"os"
	"testing"
	"time"
)

type CrudExampleTestSuite struct {
	suite.Suite
	client database.SourceInterface
}

func Test_CrudExample(t *testing.T) {
	suite.Run(t, new(CrudExampleTestSuite))
}

func (suite *CrudExampleTestSuite) SetupTest() {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	dsn := os.Getenv("MONGO_DSN")

	opts := []database.Option{
		database.Dsn(dsn),
		database.Context(ctx),
	}
	db, err := database.NewDatabase(opts...)
	if err != nil {
		assert.FailNow(suite.T(), "New Database init failed", "%v", err)
	}

	suite.client = db
}

func (suite *CrudExampleTestSuite) TearDownTest() {
	err := suite.client.Drop()
	if err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	err = suite.client.Close()
	if err != nil {
		suite.FailNow("Database close failed", "%v", err)
	}
}

func (suite *CrudExampleTestSuite) TestCrudExample_Insert() {
	err := insert(suite.client)
	assert.NoError(suite.T(), err)

	collectionMock := &mocks.CollectionInterface{}
	collectionMock.On("InsertMany", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("insert error"))

	clientMock := &mocks.SourceInterface{}
	clientMock.On("Collection", mock.Anything).Return(collectionMock, nil)

	err = insert(clientMock)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "insert error")
}

func (suite *CrudExampleTestSuite) TestCrudExample_FindMany() {
	err := insert(suite.client)
	assert.NoError(suite.T(), err)
	res, err := findMany(suite.client)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), res, 2)

	var resSlice []string
	for _, v := range res {
		resSlice = append(resSlice, v.String)
	}

	assert.Contains(suite.T(), resSlice, "value1")
	assert.Contains(suite.T(), resSlice, "value2")

	collectionMock := &mocks.CollectionInterface{}
	collectionMock.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("find error"))

	clientMock := &mocks.SourceInterface{}
	clientMock.On("Collection", mock.Anything).Return(collectionMock, nil)

	_, err = findMany(clientMock)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "find error")

	cursorMock := &mocks.CursorInterface{}
	cursorMock.On("All", mock.Anything, mock.Anything).Return(errors.New("cursor error"))

	collectionMock = &mocks.CollectionInterface{}
	collectionMock.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(cursorMock, nil)

	clientMock = &mocks.SourceInterface{}
	clientMock.On("Collection", mock.Anything).Return(collectionMock, nil)

	_, err = findMany(clientMock)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "cursor error")
}

func (suite *CrudExampleTestSuite) TestCrudExample_FindOne() {
	err := insert(suite.client)
	assert.NoError(suite.T(), err)
	res, err := findOne(suite.client)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), res.String, "value3")

	singleResultMock := &mocks.SingleResultInterface{}
	singleResultMock.On("Decode", mock.Anything).Return(errors.New("single result error"))

	collectionMock := &mocks.CollectionInterface{}
	collectionMock.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(singleResultMock, nil)

	clientMock := &mocks.SourceInterface{}
	clientMock.On("Collection", mock.Anything).Return(collectionMock, nil)

	_, err = findOne(clientMock)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "single result error")
}
