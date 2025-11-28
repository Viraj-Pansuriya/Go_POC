package model

import (
	"fmt"
	"math/rand"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Error      error
}

func FetchURL(url string, ch chan Result) Result {
	fmt.Println("Fetching Response from url", url)

	delay := rand.Intn(500)
	delayDr := time.Duration(delay) * time.Millisecond
	time.Sleep(delayDr)
	result := Result{
		URL:        url,
		StatusCode: 200,
		Duration:   delayDr,
		Error:      nil,
	}

	ch <- result
	fmt.Println("Response has been fetched for url", url)
	return result
}
