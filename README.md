PaySuper MongoDB Driver
=====

[![Build Status](https://github.com/paysuper/paysuper-database-mongo/workflows/Build/badge.svg?branch=develop)](https://github.com/paysuper/paysuper-database-mongo/actions) 
[![codecov](https://codecov.io/gh/paysuper/paysuper-database-mongo/branch/master/graph/badge.svg)](https://codecov.io/gh/paysuper/paysuper-database-mongo)
[![go report](https://goreportcard.com/badge/github.com/paysuper/paysuper-database-mongo)](https://goreportcard.com/report/github.com/paysuper/paysuper-database-mongo)

## Installation

Use go get.

	go get gopkg.in/paysuper/paysuper-database-mongo.v2

Then import the validator package into your own code.

	import "gopkg.in/paysuper/paysuper-database-mongo.v2"
	
## Usage

```go
import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    database "gopkg.in/paysuper/paysuper-database-mongo.v2"
    "log"
)

func main() {
    opts := []database.Option{
        database.Dsn("mongodb://localhost:27017/example_db"),
    }
    mongodb, err := database.NewDatabase(opts...)
    if err != nil {
        log.Fatal("MongoDB connection failed")
    }

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
    
    _, err := mongodb.Collection(collectionName).InsertMany(context.Background(), docs)
    
    if err != nil {
        log.Fatal("Data insert failed")
    }
}
```

More examples available in [examples directory](./examples)

## Environment variables:

| Name               | Required | Default  | Description                     |
|:-------------------|:--------:|:---------|:--------------------------------|
| MONGO_DIAL_TIMEOUT | -        | 10       | MongoBD dial timeout in seconds |
| MONGO_DSN          | true     | -        | MongoBD DSN connection string   |

## Contributing
We feel that a welcoming community is important and we ask that you follow PaySuper's [Open Source Code of Conduct](https://github.com/paysuper/code-of-conduct/blob/master/README.md) in all interactions with the community.

PaySuper welcomes contributions from anyone and everyone. Please refer to each project's style and contribution guidelines for submitting patches and additions. In general, we follow the "fork-and-pull" Git workflow.

The master branch of this repository contains the latest stable release of this component.

 
