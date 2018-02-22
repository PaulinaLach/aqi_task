package main

import (
	"aqi/dbImporter"
	"aqi/dbProcessor"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	importer.Starter()
	processor.Starter()
}
