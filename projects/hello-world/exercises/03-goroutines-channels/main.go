package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// Stack trace starts with "goroutine 123 [..."
	str := string(buf[:n])
	field := strings.Fields(str)[1] // Get "123"
	id, _ := strconv.Atoi(field)
	return id
}
func main() {

	ch := make(chan int) // Create channel ONCE outside loop

	// Start goroutines that SEND to channel
	for item := 0; item < 5; item++ {
		go func(i int) {
			gid := GetGoroutineID()
			fmt.Printf("Goroutine ID %d: Sending %d\n", gid, i)
			ch <- i
			fmt.Printf("Goroutine ID %d: Sent %d\n", gid, i)
		}(item)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("end")

	for i := 0; i < 5; i++ {
		result := <-ch // Main RECEIVES (blocks until data arrives)
		fmt.Println("Received:", result)
	}
	time.Sleep(1 * time.Second)
	fmt.Println("done")
}

func simulateDelay(i int) int {
	time.Sleep(1 * time.Second)
	fmt.Println("simulated delay for ", i)
	return i
}
