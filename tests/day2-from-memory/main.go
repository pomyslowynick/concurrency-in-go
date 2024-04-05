package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func main() {
	fmt.Println("Hello world")

	var lock sync.Mutex
	var wg sync.WaitGroup

	for _, filename := range os.Args[1:] {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("couldn't open the file: ", err)
			return
		}
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			lock.Lock()
			fmt.Println(scanner.Text())
			lock.Unlock()
		}
	}
	//wg.Add(1)
	go func() { fmt.Println("Hello from the thread!") }()
	wg.Wait()
}
