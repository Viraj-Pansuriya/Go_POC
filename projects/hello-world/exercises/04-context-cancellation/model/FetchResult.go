package model

import (
	"context"
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

func FetchURL(ctx context.Context, url string, ch chan Result) Result {
	//fmt.Println("Fetching Response from url", url)
	startTime := time.Now()
	select {
	case <-time.After(time.Duration(rand.Intn(100)) * time.Millisecond):
		fmt.Println("working...")
		result := Result{
			URL:        url,
			StatusCode: 200,
			Duration:   time.Since(startTime),
			Error:      nil,
		}
		ch <- result
		fmt.Println("Response has been fetched for url", url)
		return result
	case <-ctx.Done():
		fmt.Println("ctx cancelled, exiting worker for url", url)
		res := Result{url, 500, time.Since(startTime), ctx.Err()}
		ch <- res
		return res
	}
}
