package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	type value struct {
		value int
		mu    sync.Mutex
	}

	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()

		time.Sleep(2 * time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}

func mutexExample() {
	var memorySync sync.Mutex
	var data int
	go func() {
		memorySync.Lock()
		data++
		memorySync.Unlock()
	}()
	memorySync.Lock()
	if data == 0 {
		fmt.Println("data is 0")
	} else {
		fmt.Println("data is 1")
	}
}
