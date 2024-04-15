package main

import (
	"fmt"
	"math/rand"
)

func main() {

	bridge := func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-chanStream:
					if ok == false {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}
				for val := range orDone(done, stream) {
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}
	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v", v)
	}
}

func orDone(done <-chan interface{}, stream <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-stream:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func exerciseNaivePrimeCalculator() {

	randNums := func(done <-chan interface{}) <-chan interface{} {
		outStream := make(chan interface{})

		go func() {
			defer close(outStream)
			for {
				select {
				case <-done:
					return
				case outStream <- rand.Intn(20):
				}
			}
		}()
		return outStream
	}

	repeater := func(done <-chan interface{}, items ...interface{}) <-chan interface{} {
		outStream := make(chan interface{})

		go func() {
			defer close(outStream)
			for {
				for _, v := range items {
					select {
					case <-done:
						return
					case outStream <- v:
					}
				}
			}
		}()
		return outStream
	}

	take := func(done, inStream <-chan interface{}, size int) <-chan interface{} {
		outStream := make(chan interface{})

		go func() {
			defer close(outStream)
			for range size {
				select {
				case <-done:
					return
				case outStream <- <-inStream:
				}
			}
		}()
		return outStream
	}

	toInt := func(done, intStream <-chan interface{}) <-chan int {
		outStream := make(chan int)

		go func() {
			defer close(outStream)
			for v := range intStream {
				select {
				case <-done:
					return
				case outStream <- v.(int):
				}
			}
		}()
		return outStream
	}

	fnApply := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		outStream := make(chan interface{})

		go func() {
			defer close(outStream)
			for {
				select {
				case <-done:
					return
				case outStream <- fn():
				}
			}
		}()
		return outStream
	}

	primeNaive := func(done <-chan interface{}, valueStream <-chan int) <-chan interface{} {
		outStream := make(chan interface{})

		go func() {
			defer close(outStream)
			for value := range valueStream {
				fmt.Println("before decrementing", value)
				prime := true
				for divisor := value - 1; divisor > 1; divisor-- {
					if value%divisor == 0 {
						prime = false
						break
					}
				}
				if prime {
					select {
					case <-done:
						return
					case outStream <- value:
					}
				}
			}
		}()
		return outStream
	}

	done := make(chan interface{})
	defer close(done)
	for v := range take(done, randNums(done), 10) {
		fmt.Println(v)
	}
	fmt.Println("End of the rand nums\n")
	for v := range take(done, repeater(done, 1, 2, 4), 10) {
		fmt.Println(v)
	}
	fmt.Println("End of the repeater\n")

	for v := range take(done, fnApply(done, func() interface{} { return rand.Intn(200) }), 10) {
		fmt.Println(v)
	}
	fmt.Println("End of the fnApply\n")

	for v := range take(done, primeNaive(done, toInt(done, fnApply(done, func() interface{} { return rand.Intn(200) }))), 10) {
		fmt.Println(v)
	}
	fmt.Println("End of the primeNaive\n")
}
