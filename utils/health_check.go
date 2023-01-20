package utils

import (
	"fmt"
	"net/http"
	"time"
)

func GetResponse(url string) (int, int64) {
	// handle errors for no such host
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		// 0 is not a valid status code, so we use it as an error code
		PrintErr(fmt.Sprintf("unable to get response from %v: %v", url, err))
		return 0, 0
	}
	return resp.StatusCode, time.Since(start).Milliseconds()
}
