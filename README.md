# AQI

The aim of the project is to effectively process and aggregate data
to get an average aqi per country from a data file using MapReduce.

## Requirements

Make sure you are using go >= 1.9.

## Configuration

Configuration is made using env variables:
* DATA_FILE (default: "polData.json")
* MONGO_DB_NAME (default: "aq")
* MONGO_DB_COLLECTION (default: "measurements")
* MONGO_DB_URL - (default: "mongodb://localhost")
* QUERY_COUNTRY - set to query only by given country (optional)

## Running AQI

`go run main.go`

### Running with selected country

`QUERY_COUNTRY=Poland go run main.go`

## Running tests

* `cd tests`
* `go test`