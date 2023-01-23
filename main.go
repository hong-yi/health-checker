package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"log"
	"net/http"
	"os"
	"sph_assignment/handlers"
	"sph_assignment/models"
	"sph_assignment/utils"
	"sync"
	"time"
)

var wg sync.WaitGroup
var urlStatusMap = map[string]models.UrlStatus{}

var bucketName = os.Getenv("APP_BUCKET_NAME")
var inputFilename = os.Getenv("APP_INPUT_FILENAME")

func main() {

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for true {
			getApiHealthTask()
			//testApi()
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

	latencies := []types.MetricDatum{}
	responseCodes := []types.MetricDatum{}

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
	fmt.Println(len(urls))

	mutex := &sync.Mutex{}
	for _, url := range utils.GetUrlList(urlsMap) {
		go func(url string) {
			defer wg.Done()
			respCode, latency := utils.GetResponse(url)

			mutex.Lock()
			currentUrlStatus := models.UrlStatus{
				Name:       url,
				StatusCode: respCode,
				Latency:    latency,
			}
			urlsMap[url] = currentUrlStatus
			latencies = append(latencies, utils.CreateLatencyDatum(url, latency))
			responseCodes = append(responseCodes, utils.CreateStatusDatum(url, respCode))
			mutex.Unlock()
		}(url)
	}

	wg.Wait()

	utils.PutMetric(latencies)
	utils.PutMetric(responseCodes)
	utils.PrintInfo("Done! Waiting for next run...")

	// to json
	resJson, err := json.Marshal(urlsMap)
	if err != nil {
		log.Printf("error converting map to json: %v", err)
	}
	// write to s3 bucket
	publicBucketName := os.Getenv("APP_PUBLIC_BUCKET_NAME")
	utils.UploadFileToS3(resJson, publicBucketName, "latest.json")
	//fmt.Println(string(resJson))
}
