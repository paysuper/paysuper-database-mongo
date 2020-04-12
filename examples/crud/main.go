package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	database "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"log"
)

const (
	collectionName = "example_collection"
)

type Example struct {
	String string  `bson:"string"`
	Int    int     `bson:"int"`
	Float  float64 `bson:"float"`
}

func main() {
	opts := []database.Option{
		database.Dsn("mongodb://localhost:27017/example_db"),
	}
	mongodb, err := database.NewDatabase(opts...)
	if err != nil {
		log.Fatal("MongoDB connection failed")
	}

	err = insert(mongodb)
	if err != nil {
		log.Fatal("Example data insert failed")
	}

	_, err = findMany(mongodb)
	if err != nil {
		log.Fatal("find query failed")
	}

	_, err = findOne(mongodb)
	if err != nil {
		log.Fatal("findOne query failed")
	}

	log.Println("application successfully finished")
}

func insert(mongodb database.SourceInterface) error {
	ctx := context.Background()
	docs := []interface{}{
		&Example{
			String: "value1",
			Int:    1,
			Float:  11.11,
		},
		&Example{
			String: "value2",
			Int:    2,
			Float:  22.22,
		},
		&Example{
			String: "value3",
			Int:    2,
			Float:  33.33,
		},
	}

	_, err := mongodb.Collection(collectionName).InsertMany(ctx, docs)

	if err != nil {
		return err
	}

	return nil
}

func findMany(mongodb database.SourceInterface) ([]*Example, error) {
	ctx := context.Background()
	filter := bson.M{
		"string": bson.M{
			"$in": []string{"value1", "value2"},
		},
	}
	cursor, err := mongodb.Collection(collectionName).Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var findResults []*Example
	err = cursor.All(ctx, &findResults)

	if err != nil {
		return nil, err
	}

	return findResults, nil
}

func findOne(mongodb database.SourceInterface) (*Example, error) {
	var res *Example

	ctx := context.Background()
	filter := bson.M{"string": "value3"}
	err := mongodb.Collection(collectionName).FindOne(ctx, filter).Decode(&res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func transaction(mongodb database.SessionInterface) {
	/*session, err := mongodb.StartSession()

	if err != nil {
		t.Fatal(err)
	}
	if err = session.StartTransaction(); err != nil {
		t.Fatal(err)
	}*/
}
