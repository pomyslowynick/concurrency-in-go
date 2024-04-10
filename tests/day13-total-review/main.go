package main

import (
	"sync"
	"time"
)

func main() {
	condition := sync.NewCond(&sync.Mutex{})
	counter := 0

	signaller := func() {
		time.Sleep(1 * time.Second)
		condition.Signal()
	}

	go signaller()
	for {
		condition.L.Lock()
		counter++
		condition.L.Unlock()
	}
	condition.Wait()
}

func deadlockExample() {
	var mu sync.Mutex
	var mu2 sync.Mutex
	var wg sync.WaitGroup

	someFunc := func(mu, mu2 *sync.Mutex) {
		defer wg.Done()
		defer mu.Unlock()
		defer mu2.Unlock()

		mu.Lock()
		time.Sleep(1 * time.Second)
		mu2.Lock()
	}

	wg.Add(2)
	go someFunc(&mu, &mu2)
	go someFunc(&mu2, &mu)
	wg.Wait()
}
