PaySuper MongoDB Driver
=====

[![Build Status](https://travis-ci.org/paysuper/paysuper-database-mongo.svg?branch=master)](https://travis-ci.org/paysuper/paysuper-database-mongo) 
[![codecov](https://codecov.io/gh/paysuper/paysuper-database-mongo/branch/master/graph/badge.svg)](https://codecov.io/gh/paysuper/paysuper-database-mongo)
[![go report](https://goreportcard.com/badge/github.com/paysuper/paysuper-database-mongo)](https://goreportcard.com/report/github.com/paysuper/paysuper-database-mongo)

## Environment variables:

| Name               | Required | Default  | Description                                                                                                   |
|:-------------------|:--------:|:---------|:--------------------------------------------------------------------------------------------------------------|
| MONGO_HOST         | -        | -        | MongoDB host including port if this needed. **This variable is required if not set variable "MONGO_DNS"**     |
| MONGO_DB           | -        | -        | MongoDB database name. **This variable is required if not set variable "MONGO_DNS"**                          |
| MONGO_USER         | -        | ""       | MongoDB user for access to database                                                                           |
| MONGO_PASSWORD     | -        | ""       | MongoBD password for access to database                                                                       |
| MONGO_DIAL_TIMEOUT | -        | 10       | MongoBD dial timeout in seconds                                                                               |
| MONGO_DNS          | -        | -        | MongoBD DNS connection string. **This variable is required if not set variables "MONGO_HOST" and "MONGO_DB"** |
