package primes_test

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/busser/primes"
)

func TestWithSlice(t *testing.T)          { testWith(primes.WithSlice, t) }
func TestWithGoroutines(t *testing.T)     { testWith(primes.WithGoroutines, t) }
func TestWithPostponedSieve(t *testing.T) { testWith(primes.WithPostponedSieve, t) }

func testWith(f func(int, func(int)), t *testing.T) {
	testCases := []struct {
		count  int
		primes []int
	}{
		{1, []int{2}},
		{3, []int{2, 3, 5}},
		{6, []int{2, 3, 5, 7, 11, 13}},
		{10, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}},
		{9592, fixture("testdata/first-9592-primes.txt")},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d", tc.count), func(t *testing.T) {
			primes := []int{}

			f(tc.count, func(p int) {
				primes = append(primes, p)
			})

			if !reflect.DeepEqual(tc.primes, primes) {
				t.Fatalf("expected %v, got %v", tc.primes, primes)
			}
		})
	}
}

func BenchmarkWithSlice1K(b *testing.B)   { benchmarkWith(primes.WithSlice, 1024, b) }
func BenchmarkWithSlice3K(b *testing.B)   { benchmarkWith(primes.WithSlice, 3*1024, b) }
func BenchmarkWithSlice10K(b *testing.B)  { benchmarkWith(primes.WithSlice, 10*1024, b) }
func BenchmarkWithSlice30K(b *testing.B)  { benchmarkWith(primes.WithSlice, 30*1024, b) }
func BenchmarkWithSlice100K(b *testing.B) { benchmarkWith(primes.WithSlice, 100*1024, b) }
func BenchmarkWithSlice300K(b *testing.B) { benchmarkWith(primes.WithSlice, 300*1024, b) }
func BenchmarkWithSlice1M(b *testing.B)   { benchmarkWith(primes.WithSlice, 1024*1024, b) }

func BenchmarkWithGoroutines100(b *testing.B) { benchmarkWith(primes.WithGoroutines, 100, b) }
func BenchmarkWithGoroutines300(b *testing.B) { benchmarkWith(primes.WithGoroutines, 300, b) }
func BenchmarkWithGoroutines1K(b *testing.B)  { benchmarkWith(primes.WithGoroutines, 1024, b) }
func BenchmarkWithGoroutines3K(b *testing.B)  { benchmarkWith(primes.WithGoroutines, 3*1024, b) }
func BenchmarkWithGoroutines10K(b *testing.B) { benchmarkWith(primes.WithGoroutines, 10*1024, b) }

func BenchmarkWithPostponedSieve1K(b *testing.B) { benchmarkWith(primes.WithPostponedSieve, 1024, b) }
func BenchmarkWithPostponedSieve3K(b *testing.B) { benchmarkWith(primes.WithPostponedSieve, 3*1024, b) }
func BenchmarkWithPostponedSieve10K(b *testing.B) {
	benchmarkWith(primes.WithPostponedSieve, 10*1024, b)
}
func BenchmarkWithPostponedSieve30K(b *testing.B) {
	benchmarkWith(primes.WithPostponedSieve, 30*1024, b)
}
func BenchmarkWithPostponedSieve100K(b *testing.B) {
	benchmarkWith(primes.WithPostponedSieve, 100*1024, b)
}
func BenchmarkWithPostponedSieve300K(b *testing.B) {
	benchmarkWith(primes.WithPostponedSieve, 300*1024, b)
}
func BenchmarkWithPostponedSieve1M(b *testing.B) {
	benchmarkWith(primes.WithPostponedSieve, 1024*1024, b)
}

func benchmarkWith(f func(int, func(int)), count int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		f(count, func(p int) {})
	}
}

func fixture(filename string) []int {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open %s: %v", filename, err)
	}

	ints := []int{}

	s := bufio.NewScanner(f)

	for s.Scan() {
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Fatalf("failed to parse %v as int: %v", s.Text(), err)
		}
		ints = append(ints, i)
	}
	if err := s.Err(); err != nil {
		log.Fatalf("failed to read file contents: %v", err)
	}

	return ints
}
