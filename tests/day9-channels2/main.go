package main

import "fmt"

func main() {
	channelCreator := func() <-chan int {
		stream := make(chan int, 5)
		go func() {
			defer close(stream)
			for i := 0; i < 5; i++ {
				stream <- i
			}
		}()
		return stream
	}

	roChannel := channelCreator()
	for value := range roChannel {
		fmt.Printf("%v", value)
	}
}

func failingReadWrite() {
	var stream chan int
	stream <- 6
	failingInt := <-stream
	fmt.Printf("%v", failingInt)
}
