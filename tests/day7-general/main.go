package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func main() {
	//condExample()
	rwmutexExample()
}

func rwmutexExample() {
	var rwm sync.RWMutex

	producer := func(wg *sync.WaitGroup, mu sync.Locker) {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			mu.Lock()
			time.Sleep(1 * time.Second)
			mu.Unlock()
		}
	}

	observer := func(wg *sync.WaitGroup, mu sync.Locker) {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
	}

	table := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer table.Flush()
	var wg sync.WaitGroup

	test := func(count int, mu, rwm sync.Locker) time.Duration {
		startTest := time.Now()
		wg.Add(count + 1)
		go producer(&wg, mu)
		for i := count; i < 0; i-- {
			go observer(&wg, rwm)
		}
		wg.Wait()
		return time.Since(startTest)
	}

	fmt.Println("Count\tMutex\tRWMutex")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(table, "%d\t%v\t%v\n", count,
			test(count, &rwm, &rwm),
			test(count, &rwm, rwm.RLocker()),
		)
	}
}

func condExample() {
	cond := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(duration time.Duration) {
		time.Sleep(duration)

		cond.L.Lock()
		fmt.Println("Removing from the queue")
		queue = queue[1:]
		cond.L.Unlock()
		cond.Signal()
	}

	for i := 0; i < 10; i++ {
		cond.L.Lock()
		for len(queue) == 2 {
			cond.Wait()
		}
		fmt.Println("Adding to the queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		cond.L.Unlock()
	}
}
