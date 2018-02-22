package tests

import (
	"aqi/dbImporter"
	"aqi/dbProcessor"
	. "aqi/helpers"
	"gopkg.in/mgo.v2"
	"os"
	"testing"
)

func prepareDatabase(inputTestDbData []importer.DbMeasurement) (*mgo.Session, *mgo.Session) {
	dbs, err := mgo.Dial(FetchEnv("MONGO_DB_URL"))
	if err != nil {
		panic(err)
	}

	s := dbs.Clone()

	s.DB(FetchEnv("MONGO_DB_NAME")).DropDatabase()
	m := s.DB(FetchEnv("MONGO_DB_NAME")).C(FetchEnv("MONGO_DB_COLLECTION"))

	docs := make([]interface{}, len(inputTestDbData))
	for i, v := range inputTestDbData {
		docs[i] = v
	}

	err = m.Insert(docs...)
	if err != nil {
		panic(err)
	}

	return dbs, s
}

//TestProcessAqi tests if ProcessAqi runs MapReduce correctly.
func TestProcessAqi(t *testing.T) {
	testCases := []struct {
		inputTestDbData []importer.DbMeasurement
		expectedResults []processor.MapReduceResult
	}{
		{
			[]importer.DbMeasurement{
				{City: "Poland", Aqi: 90},
				{City: "Poland", Aqi: 147},
				{City: "Vietnam", Aqi: 99},
				{City: "Vietnam", Aqi: 214},
				{City: "Vietnam", Aqi: 214},
				{City: "Israel", Aqi: 74},
			},
			[]processor.MapReduceResult{
				{Id: "Israel", Value: 74},
				{Id: "Poland", Value: 118.5},
				{Id: "Vietnam", Value: 175.6666666666},
			},
		}, {
			[]importer.DbMeasurement{
				{City: "Poland", Aqi: 90},
				{City: "Vietnam", Aqi: 99},
				{City: "Israel", Aqi: 74},
			},
			[]processor.MapReduceResult{
				{Id: "Israel", Value: 74},
				{Id: "Poland", Value: 90},
				{Id: "Vietnam", Value: 99},
			},
		},
	}

	for _, testCase := range testCases {
		dbs, s := prepareDatabase(testCase.inputTestDbData)

		results, err := processor.ProcessAqi(s)
		if err != nil {
			t.Errorf("Error while processing")
		}

		for i, result := range *results {
			if result != testCase.expectedResults[i] {
				t.Errorf("Expected %+v, got %+v.", testCase.expectedResults[i], result)
			}
		}

		s.Close()
		dbs.Close()
	}
}

//TestProcessAqi tests if ProcessAqi runs MapReduce correctly
//with country to query given.
func TestProcessAqiQueryParam(t *testing.T) {
	testCase := struct {
		inputTestDbData []importer.DbMeasurement
		expectedResults []processor.MapReduceResult
	}{
		[]importer.DbMeasurement{
			{City: "Poland", Aqi: 90},
			{City: "Poland", Aqi: 147},
			{City: "Vietnam", Aqi: 99},
			{City: "Vietnam", Aqi: 214},
		},
		[]processor.MapReduceResult{
			{Id: "Vietnam", Value: 156.5},
		},
	}

	dbs, s := prepareDatabase(testCase.inputTestDbData)
	defer dbs.Close()
	defer s.Close()

	os.Setenv("QUERY_COUNTRY", "Vietnam")
	results, err := processor.ProcessAqi(s)
	if err != nil {
		t.Errorf("Error while processing")
	}
	os.Unsetenv("QUERY_COUNTRY")

	for i, result := range *results {

		if result != testCase.expectedResults[i] {
			t.Errorf("Expected %+v, got %+v.", testCase.expectedResults[i], result)
		}
	}
}
