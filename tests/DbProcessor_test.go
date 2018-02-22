package tests

import (
	"aqi/dbImporter"
	"aqi/dbProcessor"
	. "aqi/helpers"
	"fmt"
	"gopkg.in/mgo.v2"
	"testing"
)

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
		dbs, err := mgo.Dial(FetchEnv("MONGO_DB_URL"))
		if err != nil {
			panic(err)
		}

		defer dbs.Close()
		s := dbs.Clone()
		defer s.Close()

		s.DB(FetchEnv("MONGO_DB_NAME")).DropDatabase()
		m := s.DB(FetchEnv("MONGO_DB_NAME")).C(FetchEnv("MONGO_DB_COLLECTION"))

		docs := make([]interface{}, len(testCase.inputTestDbData))
		for i, v := range testCase.inputTestDbData {
			docs[i] = v
		}

		err = m.Insert(docs...)
		if err != nil {
			fmt.Println(err)
		}

		results, err := processor.ProcessAqi(s)
		for i, result := range *results {
			if result != testCase.expectedResults[i] {
				t.Errorf("Expected %+v, got %+v.", testCase.expectedResults[i], result)
			}
		}
	}
}
