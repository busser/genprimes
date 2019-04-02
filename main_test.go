package main

import (
	"fmt"
	"reflect"
	"testing"
)

var funcs = []struct {
	name string
	f    func(int, func(int))
}{
	{"slice", genPrimesWithSlice},
	{"goroutines", genPrimesWithGoroutines},
	{"postponedSieve", genPrimesWithPostponedSieve},
}

func TestGenPrimes(t *testing.T) {
	testCases := []struct {
		count  int
		primes []int
	}{
		{1, []int{2}},
		{3, []int{2, 3, 5}},
		{6, []int{2, 3, 5, 7, 11, 13}},
		{10, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}},
	}

	for _, genPrimes := range funcs {
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s/%d", genPrimes.name, tc.count), func(t *testing.T) {
				primes := []int{}

				genPrimes.f(tc.count, func(p int) {
					primes = append(primes, p)
				})

				if !reflect.DeepEqual(tc.primes, primes) {
					t.Fatalf("expected %v, got %v", tc.primes, primes)
				}
			})
		}
	}
}

func BenchmarkGenPrimes(b *testing.B) {
	for _, genPrimes := range funcs {
		for count := 1; count <= 1024; count *= 2 {
			b.Run(fmt.Sprintf("%s/%d", genPrimes.name, count), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					genPrimes.f(count, func(p int) {})
				}
			})
		}
	}
}
