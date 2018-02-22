package processor

import (
	. "aqi/helpers"
	"fmt"
	"gopkg.in/mgo.v2"
)

//Type MapReduceResult is a type of data returned by MapReduce function.
type MapReduceResult struct {
	Id    string "_id"
	Value float32
}

func mapReduceConfig() *mgo.MapReduce {
	return &mgo.MapReduce{
		Map: "function() { emit(this.city, { number: this.aqi, count: 1 }) }",
		Reduce: `
			function(key, values) {
				var a = values.shift();
				for (value of values) {
					a.number += value.number;
                    a.count += value.count;
				}
				return a;
			}
		`,
		Finalize: "function(key, value) { return value.number / value.count }",
	}
}

//ProcessAqi runs MapReduce on a collection and returns a list of average aqi for country.
func ProcessAqi(s *mgo.Session) (*[]MapReduceResult, error) {
	var result []MapReduceResult
	collection := s.DB(FetchEnv("MONGO_DB_NAME")).C(FetchEnv("MONGO_DB_COLLECTION"))
	_, err := collection.Find(nil).MapReduce(mapReduceConfig(), &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Starter() {
	dbs, err := mgo.Dial(FetchEnv("MONGO_DB_URL"))
	if err != nil {
		panic(err)
	}

	defer dbs.Close()
	s := dbs.Clone()
	defer s.Close()

	mapReduceResult, err := ProcessAqi(s)
	if err != nil {
		panic(err)
	}

	for _, item := range *mapReduceResult {
		fmt.Printf("%s - %f\n", item.Id, item.Value)
	}
}
