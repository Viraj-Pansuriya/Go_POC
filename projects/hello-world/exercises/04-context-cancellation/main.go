package main

import (
	"strconv"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/04-context-cancellation/service"
)

func main() {

	rs := &service.ResponseFetcherImpl{}

	urls := make([]string, 6)

	for index := 0; index < 6; index++ {
		urls[index] = "http://www.google.com" + strconv.FormatInt(int64(index), 10)
	}
	//rs.FetchWithContextCancel(urls)
	rs.FetchWithTimeOut(urls, 100)
}
