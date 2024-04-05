package main

import (
	"sync"
	"time"
)

type Value struct {
	mu    sync.Mutex
	value int
}

func main() {

	someFunc := func(a, b Value) {
		a.mu.Lock()
		a.value++
		a.mu.Unlock()

		time.Sleep(2 * time.Second)
		b.mu.Lock()
		b.value++
		b.mu.Unlock()

	}

	var mu sync.Mutex
	a := Value{mu: mu, value: 1}
	var mu2 sync.Mutex
	b := Value{mu: mu2, value: 1}

	var wg sync.WaitGroup
	wg.Add(2)
	go someFunc(a, b)
	go someFunc(b, a)
	wg.Wait()
}
