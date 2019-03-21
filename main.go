package main

import (
	"fmt"
)

func main() {
	genPrimesWithGoroutines(10, func(p int) { fmt.Println(p) })
}

func genPrimesWithSlice(count int, export func(int)) {
	primes := []int{}

	for n := 2; len(primes) < count; n++ {
		isPrime := true

		for _, p := range primes {
			if n%p == 0 {
				isPrime = false
				break
			}
		}

		if isPrime {
			primes = append(primes, n)
			export(n)
		}
	}
}

func genPrimesWithGoroutines(count int, export func(int)) {
	stop := make(chan bool)
	defer func() { close(stop) }()

	primes := make(chan int)

	go genInts(2, primes, stop)

	for i := 0; i < count; i++ {
		p := <-primes
		export(p)

		newPrimes := make(chan int)
		go filterMultiples(p, primes, newPrimes, stop)
		primes = newPrimes
	}
}

func genInts(start int, out chan<- int, stop <-chan bool) {
	defer func() { close(out) }()

	for n := start; ; n++ {
		select {
		case <-stop:
			return
		case out <- n:
		}
	}
}

func filterMultiples(div int, in <-chan int, out chan<- int, stop <-chan bool) {
	defer func() { close(out) }()

	for n := range in {
		if n%div != 0 {
			select {
			case <-stop:
				return
			case out <- n:
			}
		}
	}
}
