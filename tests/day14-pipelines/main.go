package main

import (
	"fmt"
	"math/rand"
)

func pipelinesExercise() {
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

func main() {
	take := func(done chan interface{}, inChan <-chan interface{}, take int) chan interface{} {
		outChan := make(chan interface{})

		go func() {
			defer close(outChan)
			for i := 0; i < take; i++ {
				select {
				case <-done:
					return
				case outChan <- <-inChan:
				}
			}
		}()
		return outChan
	}

	repeater := func(done chan interface{}, values ...int) chan interface{} {
		outChan := make(chan interface{})

		go func() {
			defer close(outChan)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case outChan <- v:
					}
				}
			}
		}()
		return outChan
	}

	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})

		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}

		}()

		return valueStream
	}

	done := make(chan interface{})
	defer close(done)
	rpt := repeater(done, 1, 2)

	taker := take(done, rpt, 5)
	fmt.Println(<-taker)

	rand := func() interface{} { return rand.Int() }

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}

}
