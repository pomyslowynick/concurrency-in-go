package main

import "fmt"

func main() {

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
