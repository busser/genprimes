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
	defer close(stop)

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

func genPrimesWithPostponedSieve(count int, export func(int)) {
	primes := make(chan int)
	stop := make(chan bool)
	defer close(stop)

	go postponedSieve(primes, stop)

	for i := 0; i < count; i++ {
		prime := <-primes
		export(prime)
	}
}

func postponedSieve(out chan<- int, stop <-chan bool) {
	out <- 2
	out <- 3
	out <- 5
	out <- 7

	sieve := newDict()

	primes := make(chan int)
	go postponedSieve(primes, stop)

	<-primes
	prime := <-primes // prime == 3
	primeSquared := prime * prime

	step := 0

	for candidate := 9; ; candidate += 2 {
		if sieve.contains(candidate) { // candidate is composite
			step = sieve.pop(candidate)
		} else if candidate < primeSquared { // candidate is prime
			select {
			case <-stop:
				return
			case out <- candidate:
			}
			continue
		} else { // candidate == primeSquared
			step = 2 * prime
			prime = <-primes
			primeSquared = prime * prime
		}
		multiple := candidate + step
		for sieve.contains(multiple) {
			multiple += step
		}
		sieve.push(multiple, step)
	}
}

type dict struct {
	m map[int]int
}

func newDict() *dict {
	return &dict{
		map[int]int{},
	}
}

func (d *dict) contains(key int) bool {
	_, ok := d.m[key]
	return ok
}

func (d *dict) pop(key int) int {
	value := d.m[key]
	delete(d.m, key)
	return value
}

func (d *dict) push(key, value int) {
	d.m[key] = value
}
