package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var lock sync.Mutex

	wordFrequencyMap := make(map[string]int)
	printFile := func(filename string, wg *sync.WaitGroup) {
		defer wg.Done()

		file, err := os.Open(filename)
		if err != nil {
			panic("cannot read the file")
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()

			lock.Lock()
			wordFrequencyMap[line]++
			lock.Unlock()
		}
	}

	for _, x := range os.Args[1:] {
		wg.Add(1)
		go printFile(x, &wg)
	}

	wg.Wait()
	fmt.Println(wordFrequencyMap)
}
