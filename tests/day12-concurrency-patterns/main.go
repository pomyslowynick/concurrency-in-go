package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, v := range integers {
				select {
				case <-done:
					return
				case intStream <- v:
				}
			}
		}()
		return intStream
	}

	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
	) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(
		done <-chan interface{},
		intStream <-chan int,
		additive int,
	) <-chan int {
		additiveStream := make(chan int)
		go func() {
			defer close(additiveStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case additiveStream <- i + additive:
				}
			}
		}()
		return additiveStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 3)

	for v := range pipeline {
		fmt.Println(v)
	}
}

func streamPipelineProcessing() {
	multiply := func(value, multiplier int) int {
		return value * multiplier
	}
	add := func(value, additive int) int {
		return value + additive
	}
	ints := []int{1, 2, 3, 4, 5}
	for _, v := range ints {
		fmt.Println(add(multiply(v, 2), 1))
	}
}

func batchPipelineProcessing() {
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4, 5}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}

func handlingErrorsInGoroutines() {
	type Result struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)

			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}

				select {
				case <-done:
					return
				case results <- result:
				}

			}
		}()
		return results
	}
	done := make(chan interface{})
	defer close(done)
	errCount := 0
	urls := []string{"a", "b", "c", "https://www.google.com", "https://baddomain"}
	returnedStream := checkStatus(done, urls...)
	for result := range returnedStream {
		if result.Error != nil {
			fmt.Println("Error returned by goroutine: ", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response %v\n", result.Response.Status)
	}
}

func unhandledErrorExample() {
	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {

				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
		}()
		return responses
	}

	done := make(chan interface{})
	defer close(done)
	checkStatus(done, "https://www.google.com", "https://baddomain")
}

func preventLeakyGoroutine() {
	chanReader := func(
		done chan interface{},
		stream chan int,
	) chan interface{} {
		terminated := make(chan interface{})

		go func() {
			defer close(terminated)
			select {
			case <-done:
				fmt.Println("Terminated")
				return
			case <-stream:
				fmt.Println("Received on stream")
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := chanReader(done, nil)
	go func() {
		time.Sleep(1 * time.Second)
		close(done)
		fmt.Println("Closed done")
	}()

	<-terminated
}
