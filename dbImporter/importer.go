package importer

import (
	. "aqi/helpers"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

//Type Measurement is type of data got from input json.
type Measurement struct {
	City string
	Aqi  string
}

//Type DbMeasurement is type of data stored in database.
type DbMeasurement struct {
	City string
	Aqi  uint
}

//BatchSize is a maximum size of batch inserted to database.
/*RegNoBracketsComp is a regular expression used to remove explanations
in brackets after the country name in record of type Measurement.*/
const (
	BatchSize         = 50
	RegNoBracketsComp = `\s*\(.*\)\s*`
)

var RegNoBrackets = regexp.MustCompile(RegNoBracketsComp)

//FileReader opens the input file and sends records to a channel.
func FileReader(recordsChannel chan Measurement) {
	jsonFile, err := os.Open(FetchEnv("DATA_FILE"))

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	dec := json.NewDecoder(jsonFile)

	_, err = dec.Token()

	if err != nil {
		log.Fatal(err)
	}

	for dec.More() {
		var record Measurement
		err := dec.Decode(&record)

		if err != nil {
			log.Fatal(err)
		}

		recordsChannel <- record
	}

	close(recordsChannel)

}

func batchManager(recordsChannel chan Measurement, dbs *mgo.Session) {
	var ws sync.WaitGroup
	var records []interface{}
	var i int

	for {
		if records == nil {
			i = 0
			records = make([]interface{}, BatchSize)
		}
		record, ok := <-recordsChannel
		if !ok {
			ws.Add(1)
			go dbImporter(records[:i], dbs, &ws)
			break
		}

		transformedRecord := RecordTransformer(&record)
		if transformedRecord != nil {
			records[i] = transformedRecord
			i += 1
		}

		if i >= BatchSize-1 {
			ws.Add(1)
			go dbImporter(records[:], dbs, &ws)
			records = nil
		}
	}

	ws.Wait()
}

func dbImporter(records []interface{}, dbs *mgo.Session, ws *sync.WaitGroup) {
	defer ws.Done()

	s := dbs.Clone()
	defer s.Close()

	m := s.DB(FetchEnv("MONGO_DB_NAME")).C(FetchEnv("MONGO_DB_COLLECTION"))

	bulk := m.Bulk()
	bulk.Unordered()
	bulk.Insert(records...)
	_, bulkErr := bulk.Run()
	if bulkErr != nil {
		panic(bulkErr)
	}
}

//RecordTransformer transforms record of type Measurement to a record of type DbMeasurement.
func RecordTransformer(record *Measurement) *DbMeasurement {
	processedAqi, err := strconv.Atoi(record.Aqi)
	if err != nil {
		return nil
	}

	cityNoBrackets := RegNoBrackets.ReplaceAllString(record.City, "")
	splittedCity := strings.Split(cityNoBrackets, ", ")
	processedCity := splittedCity[len(splittedCity)-1]

	return &DbMeasurement{Aqi: uint(processedAqi), City: processedCity}
}

func Starter() {
	dbs, err := mgo.Dial(FetchEnv("MONGO_DB_URL"))
	if err != nil {
		panic(err)
	}
	defer dbs.Close()

	recordsChannel := make(chan Measurement, BatchSize)

	go FileReader(recordsChannel)

	batchManager(recordsChannel, dbs)
}
