package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

type value struct{}

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	queue := make([]int, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		cond.L.Lock()
		queue = queue[1:]
		fmt.Println("Remove from queue")
		cond.L.Unlock()
		cond.Signal()
	}

	for i := 0; i < 10; i++ {
		cond.L.Lock()
		for len(queue) == 2 {
			cond.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, 1)
		go removeFromQueue(1 * time.Millisecond)
		cond.L.Unlock()
	}
}

func rwmutex_test() {
	producer := func(wg *sync.WaitGroup, mu sync.Locker) {
		defer wg.Done()
		for i := 0; i < 6; i++ {
			mu.Lock()
			time.Sleep(1)
			mu.Unlock()
		}
	}

	observer := func(wg *sync.WaitGroup, mu sync.Locker) {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
	}

	test := func(count int, mu, rwm sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)

		startTest := time.Now()
		go producer(&wg, mu)
		for ; count > 0; count-- {
			go observer(&wg, rwm)
		}
		wg.Wait()
		return time.Since(startTest)
	}

	var mu sync.RWMutex
	table := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer table.Flush()

	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))

		fmt.Fprintf(
			table,
			"%d\t%v\t%v\n",
			count,
			test(count, &mu, mu.RLocker()),
			test(count, &mu, &mu),
		)
	}
}
