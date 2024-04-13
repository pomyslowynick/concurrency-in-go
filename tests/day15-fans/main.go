package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	tee := func(
		done <-chan interface{},
		in <-chan interface{},
	) (<-chan interface{}, <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDoneExample(done, in) {
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}

				}
			}
		}()
		return out1, out2
	}

	simpleChan := func(done <-chan interface{}) <-chan interface{} {
		outChan := make(chan interface{})

		go func() {
			defer close(outChan)

			for {
				select {
				case <-done:
					return
				case outChan <- rand.Intn(200):
				}
			}
		}()

		return outChan
	}

	done := make(chan interface{})

	outChan1, outChan2 := tee(done, simpleChan(done))

	for i := 0; i < 5; i++ {
		fmt.Println(<-outChan1)
		fmt.Println(<-outChan2)
	}
}

func orDoneExample(done, stream <-chan interface{}) <-chan interface{} {

	orDone := func(done, stream <-chan interface{}) <-chan interface{} {
		outStream := make(chan interface{})

		go func() {
			defer close(outStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-stream:
					if ok == false {
						return
					}
					select {
					case outStream <- v:
					case <-done:
					}
				}
			}
		}()
		return outStream

	}
	return orDone(done, stream)
}

func fanInOut() {

	take := func(done <-chan interface{}, valueStream <-chan interface{}, items int) <-chan interface{} {
		outStream := make(chan interface{})
		go func() {

			defer close(outStream)
			for range items {
				select {
				case <-done:
					return
				case outStream <- <-valueStream:
				}
			}
		}()
		return outStream
	}

	toInt := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()
		return intStream
	}

	primeFinder := func(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
		primeStream := make(chan interface{})
		go func() {
			defer close(primeStream)
			for integer := range intStream {
				integer -= 1
				prime := true
				for divisor := integer - 1; divisor > 1; divisor-- {
					if integer%divisor == 0 {
						prime = false
						break
					}
				}
				if prime {
					select {
					case <-done:
						return
					case primeStream <- integer:
					}
				}
			}
		}()
		return primeStream
	}

	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
	) <-chan interface{} {

		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		multiplex := func(c <-chan interface{}) {
			defer wg.Done()

			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, channel := range channels {
			go multiplex(channel)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}
	rand := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	// for prime := range take(done, primeFinder(done, randIntStream), 10) {
	// 	fmt.Printf("\t%d\n", prime)
	// }

	numFinders := runtime.NumCPU()
	finders := make([]<-chan interface{}, numFinders)

	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Println(prime)
	}
	fmt.Printf("Search took: %v", time.Since(start))

}

func pipelineFunctorExercise() {

	// TODO: Try out using generics instead of empty interfaces
	repeater := func(done chan interface{}, vals ...interface{}) chan interface{} {
		outChan := make(chan interface{})

		go func() {
			defer close(outChan)
			for {
				for _, val := range vals {
					select {
					case <-done:
						return
					case outChan <- val:
					}
				}
			}
		}()
		return outChan
	}

	done := make(chan interface{})
	repeat := repeater(done, 1, 2, 4, 5)

	functer := repeatFn(done, func() interface{} {
		return 2 * 2
	})

	for range 10 {
		fmt.Println(<-repeat)
	}
	for range 10 {
		fmt.Println(<-functer)
	}
}

func repeatFn(done chan interface{}, fn func() interface{}) chan interface{} {
	outChan := make(chan interface{})

	go func() {
		defer close(outChan)
		for {
			select {
			case <-done:
				return
			case outChan <- fn():
			}
		}
	}()
	return outChan
}
