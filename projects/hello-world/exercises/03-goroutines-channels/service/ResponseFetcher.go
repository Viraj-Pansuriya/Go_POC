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
	Fetch(urls []string)
}

type UrlResponseFetcher struct {
}

func (ur *UrlResponseFetcher) Fetch(urls []string) {
	results := make(map[string]model.Result, len(urls))

	ch := make(chan model.Result)
	startTime := time.Now()
	for index := 0; index < len(urls); index++ {
		go model.FetchURL(urls[index], ch)
	}
	for index := 0; index < len(urls); index++ {
		res := <-ch
		results[res.URL] = res
		fmt.Println("Fetched result from URL", res.URL, " is : ", res)
	}
	endTime := time.Since(startTime)
	fmt.Println("Total Time taken for batch request is ", endTime.Milliseconds())
	return
}

func GetResponseFetcherInstance() ResponseFetcher {
	if responseFetcherInstance == nil {
		responseFetcherInstance = &UrlResponseFetcher{}
	}
	return responseFetcherInstance
}
