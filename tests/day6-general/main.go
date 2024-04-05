package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(duration time.Duration) {
		time.Sleep(duration)
		cond.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed item from the queue")
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
