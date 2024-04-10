package main

import "fmt"

func main() {
	generator := func(done <-chan interface{}, ints ...int) chan int {
		outChan := make(chan int)
		go func() {
			defer close(outChan)
			for _, val := range ints {

				select {
				case <-done:
					return
				case outChan <- val:
				}
			}
		}()
		return outChan
	}

	multiplier := func(done <-chan interface{}, multiplier int, generator <-chan int) chan int {
		multiChan := make(chan int)
		go func() {
			defer close(multiChan)
			for val := range generator {
				select {
				case <-done:
					return
				case multiChan <- val * multiplier:
				}
			}

		}()
		return multiChan
	}

	adder := func(done <-chan interface{}, additive int, generator <-chan int) chan int {
		addChan := make(chan int)
		go func() {
			defer close(addChan)
			for val := range generator {
				select {
				case <-done:
					return
				case addChan <- val + additive:
				}
			}

		}()
		return addChan
	}

	done := make(chan interface{})
	mygen := adder(done, 15, multiplier(done, 2, multiplier(done, 10, generator(done, 1, 2, 4, 5))))
	defer close(done)
	for val := range mygen {
		fmt.Println(val)
	}
}
