package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sph_assignment/handlers"
	"sph_assignment/utils"
	"sync"
	"time"
)

var wg sync.WaitGroup

var bucketName = os.Getenv("APP_BUCKET_NAME")
var inputFilename = os.Getenv("APP_INPUT_FILENAME")

func main() {

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for true {
			getApiHealthTask()
			time.Sleep(time.Minute * 10)
		}
		wg.Done()
	}()

	go func() {
		// basic http handler to retrieve file from S3
		http.HandleFunc("/ping", handlers.AppHealthHandler)
		http.HandleFunc("/healthcheck", handlers.CheckHealthHandler)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			return
		}
	}()

	wg.Wait()

}

func getApiHealthTask() {
	inputFile, err := utils.DownloadFileFromS3(bucketName, inputFilename)
	if err != nil {
		utils.PrintErr(fmt.Sprintf("unable to download file %v", err))
		return
	}

	urls, err := utils.ParseFile(inputFile)
	if err != nil {
		log.Println("error")
	}
	urlsMap := utils.GetUrls(urls)
	wg.Add(len(urls))

	for url, urlStatus := range urlsMap {
		//log.Printf("%v: %v\n", url, getResponse(url.Url))
		urlStatus := urlStatus
		url := url
		go func() {
			defer wg.Done()
			urlStatus.StatusCode, urlStatus.Latency = utils.GetResponse(url)
			utils.PutStatusMetrics(url, urlStatus.StatusCode)
			if urlStatus.StatusCode == http.StatusOK {
				utils.PutLatencyMetrics(url, urlStatus.Latency)
			}
			urlsMap[url] = urlStatus
		}()
	}

	wg.Wait()

	//fmt.Println(urlsMap)
	// to json
	resJson, err := json.Marshal(urlsMap)
	if err != nil {
		log.Printf("error converting map to json: %v", err)
	}
	// write to s3 bucket
	utils.UploadFileToS3(resJson, bucketName, "latest.json")
	//fmt.Println(string(resJson))
}
