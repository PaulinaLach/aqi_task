package tests

import (
	"aqi/dbImporter"
	"os"
	"path"
	"testing"
)

//TestFileReader tests if FileReader works correctly with normal and empty file.
func TestFileReader(t *testing.T) {
	files := []struct {
		filePath       string
		processedFiles []importer.Measurement
	}{
		{
			path.Join("fixtures", "testCorrectJson.json"),
			[]importer.Measurement{
				{City: "airp. Lafarge, Poland", Aqi: "90"},
				{City: "Ho Chi Minh City US Consulate, Vietnam (Hà Nội Đại sứ quán Mỹ, Vietnam)", Aqi: "99"},
			},
		},
		{
			path.Join("fixtures", "testEmptyJson.json"),
			[]importer.Measurement{},
		},
	}

	for _, file := range files {
		os.Setenv("DATA_FILE", file.filePath)
		c := make(chan importer.Measurement)

		go importer.FileReader(c)

		i := 0
		for v := range c {
			if v != file.processedFiles[i] {
				t.Errorf("Expected %+v, got %+v", file.processedFiles[i], v)
			}
			i++
		}
	}
}

//TestRecordsTransformer tests if RecordTransformer transforms correctly
// records of type Measurement to DbMeasurement type.
func TestRecordTransformer(t *testing.T) {
	measurements := []struct {
		testMeasurement   importer.Measurement
		testDbMeasurement importer.DbMeasurement
	}{
		{importer.Measurement{City: "airp. Lafarge, Poland", Aqi: "90"},
			importer.DbMeasurement{City: "Poland", Aqi: 90}},
		{importer.Measurement{City: "Ho Chi Minh City US Consulate, Vietnam (Hà Nội Đại sứ quán Mỹ, Vietnam)", Aqi: "99"},
			importer.DbMeasurement{City: "Vietnam", Aqi: 99}},
		{importer.Measurement{City: "South Ashkelon, southern coastal plain, Israel (ישראל,אשקלון דרום, מישור החוף הדרומי)", Aqi: "74"},
			importer.DbMeasurement{City: "Israel", Aqi: 74}},
	}

	for _, measurement := range measurements {

		dbMeasurement := importer.RecordTransformer(&measurement.testMeasurement)
		if *dbMeasurement != measurement.testDbMeasurement {
			t.Errorf("Expected (%s, %u), got (%s, %u).",
				measurement.testDbMeasurement.City, measurement.testDbMeasurement.Aqi,
				dbMeasurement.City, dbMeasurement.Aqi)
		}
	}
}
