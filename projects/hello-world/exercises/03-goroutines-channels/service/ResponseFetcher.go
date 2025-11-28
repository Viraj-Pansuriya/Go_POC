package service

import (
	"fmt"
	"time"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/03-goroutines-channels/model"
)

var (
	responseFetcherInstance ResponseFetcher
)

type ResponseFetcher interface {
	Fetch(urls []string) []model.Result
}

type UrlResponseFetcher struct {
}

func (ur *UrlResponseFetcher) Fetch(urls []string) []model.Result {
	results := make([]model.Result, len(urls))
	ch := make(chan model.Result)
	startTime := time.Now()
	for index := 0; index < len(urls); index++ {
		go model.FetchURL(urls[index], ch)
	}
	for index := 0; index < len(urls); index++ {
		results[index] = <-ch
		fmt.Println("Fetched result from URL", urls[index], " is : ", results[index])
	}
	endTime := time.Since(startTime)
	fmt.Println("Total Time taken for batch request is ", endTime.Milliseconds())
	return results
}

func GetResponseFetcherInstance() ResponseFetcher {
	if responseFetcherInstance == nil {
		responseFetcherInstance = &UrlResponseFetcher{}
	}
	return responseFetcherInstance
}
