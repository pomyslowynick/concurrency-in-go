package main

import (
	"fmt"
)

func main() {
	stream := make(chan interface{})
	close(stream)
	stream2 := make(chan interface{})
	close(stream2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {

		select {
		case <-stream:
			c1Count++
		case <-stream2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}
func channelOwnerEx() {

	channelCreator := func() <-chan int {
		stream := make(chan int, 5)
		go func() {
			defer close(stream)
			for i := 1; i < 5; i++ {
				stream <- i
			}
		}()
		return stream
	}

	for value := range channelCreator() {
		fmt.Printf("%v\n", value)
	}

}
