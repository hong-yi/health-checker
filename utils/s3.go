package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

func DownloadFileFromS3(bucketName string, fileName string) ([]byte, error) {
	client := s3.NewFromConfig(GetAwsCredentials())
	fileData, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		//log.Printf("[ERROR] unable to download file from s3 %v", err)
		return nil, err
	}
	body, err := io.ReadAll(fileData.Body)
	if err != nil {
		//log.Printf("[ERROR] unable to read downloaded file %v", err)
		return nil, err
	}
	return body, nil
}

func UploadFileToS3(fileData []byte, bucketName string, fileName string) {
	client := s3.NewFromConfig(GetAwsCredentials())
	bytesReader := bytes.NewReader(fileData)
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytesReader,
	})
	if err != nil {
		PrintErr(fmt.Sprintf("unable to read filedata %v", err))
		return
	}
	PrintInfo(fmt.Sprintf("File upload successful!"))
}
