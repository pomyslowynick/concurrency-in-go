package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()

	// message := "hello"
	// var wg sync.WaitGroup
	// wg.Add(1)
	//
	//	go func() {
	//		defer wg.Done()
	//		message = "welcome"
	//	}()
	//
	// wg.Wait()
	// fmt.Println(message)
}
