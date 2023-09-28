package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	applicationsgenerator "emcontroller/auto-schedule/experiments/applications-generator"
	"emcontroller/models"
)

var (
	appNamePrefix string = "expt-app"
	appCount      int    = 100

	algorithms  []string = []string{"BERand", "Amaga", "Ampga", "Diktyoga", "Mcssga"}
	repeatCount int      = 50

	dataPath     string = "executor-python/data"
	jsonFileName string = "request_applications.json"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	for i := 1; i <= repeatCount; i++ {
		repeatPath := filepath.Join(dataPath, fmt.Sprintf("repeat%d", i))

		// create the repeat path
		log.Printf("create dir %s", repeatPath)
		if err := os.MkdirAll(repeatPath, fs.ModePerm); err != nil {
			log.Panicf("create path %s, error: %s", repeatPath, err.Error())
		}

		for _, algoName := range algorithms {
			algoPath := filepath.Join(repeatPath, algoName)

			// create the algorithm path
			log.Printf("create dir %s", algoPath)
			if err := os.MkdirAll(algoPath, fs.ModePerm); err != nil {
				log.Panicf("create path %s, error: %s", algoPath, err.Error())
			}
		}

		// create the json file to save the request body to deploy applications
		apps, err := applicationsgenerator.MakeExperimentApps(appNamePrefix, appCount, true)
		if err != nil {
			log.Panicf("MakeExperimentApps error: %s", err.Error())
		}

		func() { // to make the "defer" effective, we use this anonymous function
			jsonFilePath := filepath.Join(repeatPath, jsonFileName)
			log.Printf("write app deploy request body json in to file %s", jsonFilePath)
			jsonFile, err := os.Create(jsonFilePath)
			defer jsonFile.Close()
			if err != nil {
				log.Panicf("create file %s, error: %s", jsonFilePath, err.Error())
			}
			if _, err := jsonFile.WriteString(models.JsonString(apps)); err != nil {
				log.Panicf("Write json to file %s, error: %s", jsonFilePath, err.Error())
			}
		}()
	}

}
