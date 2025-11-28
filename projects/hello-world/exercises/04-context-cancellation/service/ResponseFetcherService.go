package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/04-context-cancellation/model"
)

type ResponseFetcher interface {
	FetchResp(url string)
}

type ResponseFetcherImpl struct{}

func (rf *ResponseFetcherImpl) FetchResp(url string) {

	// do some work here;
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan model.Result)
	go model.FetchURL(ctx, url, ch)
	resp := <-ch
	fmt.Println(resp)
	time.Sleep(2 * time.Second)
	defer cancel()

}

func (rf *ResponseFetcherImpl) SimulateContext(ctx context.Context) {

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Timeout")
	case <-ctx.Done():
		fmt.Println("Done")
	}

	fmt.Println("Outside a select statement")
}

func (rf *ResponseFetcherImpl) FetchWithContextCancel(urls []string) {
	ch := make(chan model.Result)
	for index := 0; index < len(urls); index++ {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		go model.FetchURL(ctx, urls[index], ch)
		delay := time.Duration(rand.Intn(200)) * time.Millisecond
		fmt.Println("Simulating delay for index , ", index, " is : ", delay)
		time.Sleep(delay)
		cancel()
	}

	for index := 0; index < len(urls); index++ {
		res := <-ch
		fmt.Println("Received response", res)
	}
}

func (rf *ResponseFetcherImpl) FetchWithTimeOut(urls []string, timeOut int) {

	ch := make(chan model.Result)
	for index := 0; index < len(urls); index++ {
		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, time.Duration(timeOut)*time.Millisecond)
		go model.FetchURL(ctx, urls[index], ch)
	}

	for index := 0; index < len(urls); index++ {
		res := <-ch
		fmt.Println("Received response", res)
	}

}
