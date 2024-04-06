package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}

			}
		}()
		return orDone
	}
}
func doneWriteChannel() {
	newRandStream := func(done <-chan interface{}) (<-chan int, <-chan int) {
		randStream := make(chan int)
		terminated := make(chan int)

		go func() {
			defer fmt.Println("Goroutine is closed")
			defer close(randStream)
			defer close(terminated)
			randStream <- rand.Int()
			for {
				select {
				case <-done:
					return
				case randStream <- rand.Int():
					fmt.Println("someting")
				}
			}
		}()
		return randStream, terminated
	}
	done := make(chan interface{})
	randStream, terminated := newRandStream(done)
	fmt.Println("3 random ints:")

	for i := 0; i < 3; i++ {
		fmt.Printf("%v\n", <-randStream)
	}
	go func() {
		time.Sleep(1 * time.Second)
		close(done)
	}()
	<-terminated
	fmt.Println("Done.")
}

func leakingWriteBlockedGoroutine() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("Goroutine is closed")
			defer close(randStream)
			for {
				randStream <- rand.Int()
			}
		}()
		return randStream
	}
	randStream := newRandStream()
	fmt.Println("3 random ints:")

	for i := 0; i < 3; i++ {
		fmt.Printf("%v\n", <-randStream)
	}
}

func donePatternGoroutineCancelling() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminated)

			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Cancelling doWork gorutine")
		close(done)
	}()

	// do some work

	<-terminated
	fmt.Println("Done.")
}

func leakyGoroutine() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				fmt.Println(s)

			}
		}()
		return completed
	}
	doWork(nil)
	fmt.Println("Done.")
}

func confinedBytesBuffer() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()
}

func confinedChannel() {
	data := make([]int, 4)

	data[0] = 4
	data[1] = 5
	data[2] = 6

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

func channelOwnerPattern() {
	channelCreator := func() <-chan int {
		stream := make(chan int)

		go func() {
			defer close(stream)
			for i := 0; i < 5; i++ {
				stream <- i
			}
		}()

		return stream
	}

	for variable := range channelCreator() {
		fmt.Printf("%v\n", variable)
	}
}

func selectTest() {
	stream := make(chan int)
	stream2 := make(chan int)
	stream3 := make(chan int)
	defer close(stream)
	defer close(stream2)
	defer close(stream3)

	go func() {
		stream <- 1
	}()
	go func() {
		stream2 <- 1
	}()
	go func() {
		stream3 <- 1
	}()
	for {
		select {
		case <-stream:
			fmt.Println("First stream")
		case <-stream2:
			fmt.Println("Second stream")
		case <-stream3:
			fmt.Println("Third stream")
		default:
			fmt.Println("default case")
		}
		time.Sleep(1 * time.Second)
	}
}
