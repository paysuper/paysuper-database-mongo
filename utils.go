package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CodeDuplicateKeyErrorCollection = 11000
)

func ToSortOption(fields []string) interface{} {
	sort := make(map[string]interface{})

	for _, field := range fields {
		order := 1

		if field == "" {
			continue
		}

		switch field[0] {
		case '+':
			field = field[1:]
		case '-':
			order = -1
			field = field[1:]
		}

		sort[field] = order
	}

	if len(sort) <= 0 {
		sort["_id"] = 1
	}

	return sort
}

func IsDuplicate(err error) bool {
	writeErr, ok := err.(mongo.WriteException)

	if !ok {
		return false
	}

	return writeErr.WriteErrors[0].Code == CodeDuplicateKeyErrorCollection
}

func GetObjectIDCounter(objectID primitive.ObjectID) int64 {
	b := []byte(objectID.Hex()[9:12])
	return int64(uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]))
}
