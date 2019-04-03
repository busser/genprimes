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

// solution based on:
// https://stackoverflow.com/a/12563800
// https://stackoverflow.com/a/10733621
func postponedSieve(out chan<- int, stop <-chan bool) {
	defer close(out)

	for _, prime := range []int{2, 3, 5, 7} {
		select {
		case <-stop:
			return
		case out <- prime:
		}
	}

	primes := make(chan int)
	go postponedSieve(primes, stop)

	select {
	case <-stop:
		return
	case <-primes:
	}

	var prime int
	select {
	case <-stop:
		return
	case prime = <-primes: // prime == 3
	}

	primeSquared := prime * prime

	sieve := newDict()
	var step int

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
			select {
			case <-stop:
				return
			case prime = <-primes:
			}
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
