package utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"time"
)

// insert functions to write metrics to cloudwatch
func CreateLatencyDatum(url string, latency int64) types.MetricDatum {
	metricToPut := types.MetricDatum{
		MetricName: aws.String("website_latency"),
		Dimensions: []types.Dimension{
			types.Dimension{
				Name:  aws.String("url"),
				Value: aws.String(url),
			},
		},
		Timestamp: aws.Time(time.Now()),
		Unit:      types.StandardUnitMilliseconds,
		Value:     aws.Float64(float64(latency)),
	}
	return metricToPut
}

//func CreateStatusDatum(url string, statusCode int) types.MetricDatum {
//	metricValue := 1
//	if statusCode != http.StatusOK {
//		metricValue = 0
//	}
//	fmt.Println(metricValue)
//	metricToPut := types.MetricDatum{
//		MetricName: aws.String("website_status"),
//		Dimensions: []types.Dimension{
//			types.Dimension{
//				Name:  aws.String("url"),
//				Value: aws.String(url),
//			},
//		},
//		Timestamp: aws.Time(time.Now()),
//		Unit:      types.StandardUnitCount,
//		Value:     aws.Float64(float64(metricValue)),
//	}
//	return metricToPut
//}

func PutMetric(datums []types.MetricDatum) {
	client := cloudwatch.NewFromConfig(GetAwsCredentials())
	_, err := client.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		MetricData: datums,
		Namespace:  aws.String("health_stats"),
	})
	if err != nil {
		PrintErr(fmt.Sprintf("unable to put metric data for %v", err))
		return
	}
	//PrintInfo(fmt.Sprintf("%v metric added!", *datum.MetricName))
}
