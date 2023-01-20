package utils

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
	"sph_assignment/models"
	"strings"
)

func ParseFile(fileData []byte) ([][]string, error) {
	byteReader := bytes.NewReader(fileData)
	csvReader := csv.NewReader(byteReader)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file %v", err)
	}
	// remove the headers
	return data[1:], nil
}

func GetUrls(urls [][]string) map[string]models.UrlStatus {
	urlList := map[string]models.UrlStatus{}
	for _, url := range urls {
		urlList[strings.TrimSpace(url[1])] = models.UrlStatus{Name: strings.TrimSpace(url[0])}
	}
	return urlList
}

func GetAwsCredentials() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		PrintErr(fmt.Sprintf("unable to load default config %v", err))
	}
	return cfg
}

func PrintErr(msg string) {
	log.Printf("[ERROR] %v", msg)
}

func PrintInfo(msg string) {
	log.Printf("[INFO] %v", msg)
}
