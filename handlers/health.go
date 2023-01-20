package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func CheckHealthHandler(w http.ResponseWriter, r *http.Request) {
	publicBucketName := os.Getenv("APP_PUBLIC_BUCKET_NAME")
	if len(publicBucketName) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("https://%s.s3.ap-southeast-1.amazonaws.com/latest.json", publicBucketName), http.StatusPermanentRedirect)
}

func AppHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
