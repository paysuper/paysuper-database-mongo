package database

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestToSortOption(t *testing.T) {
	sort := ToSortOption([]string{"created_at", "-id", "", "+field1"})
	assert.Len(t, sort, 3)
	assert.Contains(t, sort, "created_at")
	assert.Contains(t, sort, "id")

	tSort, ok := sort.(map[string]interface{})
	assert.True(t, ok)
	assert.EqualValues(t, tSort["created_at"], 1)
	assert.EqualValues(t, tSort["id"], -1)
	assert.EqualValues(t, tSort["field1"], 1)
}

func TestToSortOptionDefaultSort(t *testing.T) {
	sort := ToSortOption([]string{})
	assert.Len(t, sort, 1)
	assert.Contains(t, sort, "_id")
	tSort, ok := sort.(map[string]interface{})
	assert.True(t, ok)
	assert.EqualValues(t, tSort["_id"], 1)
}

func TestIsDuplicate(t *testing.T) {
	err := mongo.WriteException{
		WriteErrors: []mongo.WriteError{
			{Code: CodeDuplicateKeyErrorCollection},
		},
	}
	assert.True(t, IsDuplicate(err))
	assert.False(t, IsDuplicate(ErrorSessionNotInit))
}

func TestGetObjectIDCounter(t *testing.T) {
	counter := GetObjectIDCounter(primitive.NewObjectID())
	assert.NotZero(t, counter)
}
