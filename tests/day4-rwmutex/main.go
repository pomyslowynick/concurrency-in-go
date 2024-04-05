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
	producer := func(wg *sync.WaitGroup, mu sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			mu.Lock()
			time.Sleep(1 * time.Second)
			mu.Unlock()
		}
	}

	consumer := func(wg *sync.WaitGroup, mu sync.Locker) {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
	}
	var mu sync.RWMutex

	table := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer table.Flush()
	fmt.Fprintf(table, "Readers\tTWMutex\tMutex\n")

	test := func(mutex, rwMutex sync.Locker, count int) time.Duration {

		var wg sync.WaitGroup
		wg.Add(count + 1)
		startTime := time.Now()

		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go consumer(&wg, rwMutex)
		}
		wg.Wait()
		return time.Since(startTime)
	}

	runner := func() {
		for i := 0; i < 20; i++ {
			count := int(math.Pow(2, float64(i)))
			fmt.Fprintf(
				table,
				"%d\t%v\t%v\n",
				count,
				test(&mu, &mu, count),
				test(&mu, mu.RLocker(), count),
			)
		}
	}

	runner()
}
