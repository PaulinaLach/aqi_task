# AQI

The aim of the project is to effectively process and aggregate data
to get an average aqi per country from a data file using MapReduce.

## Requirements

Make sure you are using go >= 1.9.

## Configuration

Configuration is made using env variables:
* DATA_FILE - by default it is "polData.json"
* MONGO_DB_NAME - set to "aq"
* MONGO_DB_COLLECTION - set to "measurements"
* MONGO_DB_URL - set to "mongodb://localhost"

## Running AQI

`go run main.go`

## Running tests

* `cd tests`
* `go test`