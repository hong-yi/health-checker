package main

import (
	"encoding/json"
	"fmt"
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
	//waitChan := make(chan struct{}, MAX_CONCURRENT_JOBS)
	//count := 0
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
	wg.Add(len(urls) - 1)
	fmt.Println(len(urls))

	mutex := &sync.Mutex{}
	for _, url := range utils.GetUrlList(urlsMap) {
		//log.Printf("%v: %v\n", url, getResponse(url.Url))
		go func(url string) {
			defer wg.Done()

			//urlStatus.StatusCode, urlStatus.Latency = utils.GetResponse(url)
			//count++
			respCode, latency := utils.GetResponse(url)

			mutex.Lock()
			currentUrlStatus := models.UrlStatus{
				Name:       url,
				StatusCode: respCode,
				Latency:    latency,
			}
			urlsMap[url] = currentUrlStatus
			//results[url] = resp.StatusCode
			mutex.Unlock()
			utils.PutStatusMetrics(url, currentUrlStatus.StatusCode)
			if currentUrlStatus.StatusCode == http.StatusOK {
				utils.PutLatencyMetrics(url, currentUrlStatus.Latency)
			}
			//urlsMap[url] = urlStatus
		}(url)
	}

	wg.Wait()

	//fmt.Println(urlsMap)
	fmt.Println("done")
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
