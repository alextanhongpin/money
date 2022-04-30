package money_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/alextanhongpin/money"
)

func ExampleSplit() {
	m := money.New(5030, 1)
	s := m.Split(3)

	fmt.Println(s, money.Sum(s...))
	// Output: [1677 1677 1676] 5030
}

func ExampleAllocate() {
	m := money.New(5030, 1)
	a := m.Allocate(1, 2, 5)

	fmt.Println(a, money.Sum(a...))
	// Output: [629 1258 3143] 5030
}

func ExampleAllocateMapIntKey() {
	m := money.New(5030, 1)
	a := money.AllocateMap(m, map[int64]uint{
		1000: 1,
		2000: 2,
		3000: 5,
	})

	values := make([]int64, 0, len(a))
	for _, val := range a {
		values = append(values, val)
	}

	fmt.Println(a, money.Sum(values...))
	// Output: map[1000:629 2000:1258 3000:3143] 5030
}

func ExampleAllocateMapStrKey() {
	m := money.New(5030, 1)
	a := money.AllocateMap(m, map[string]uint{
		"a": 5,
		"b": 2,
		"c": 1,
	})

	values := make([]int64, 0, len(a))
	for _, val := range a {
		values = append(values, val)
	}

	fmt.Println(a, money.Sum(values...))
	// Output: map[a:3144 b:1258 c:628] 5030
}

func FuzzMoneySplit(f *testing.F) {
	f.Fuzz(func(t *testing.T, positive, n uint) {
		if positive < n {
			positive, n = n, positive
		}
		positive++
		n++

		amount := int64(positive)
		m := money.New(amount, 1)
		res := m.Split(n)
		if sum := money.Sum(res...); amount != sum {

			t.Errorf("split %d by %d, expected %d, got %d", amount, n, amount, sum)
		}
	})
}

func FuzzMoneyAllocate(f *testing.F) {
	f.Fuzz(func(t *testing.T, positive, n uint) {
		if positive < n {
			positive, n = n, positive
		}
		positive++
		n++

		amount := int64(positive)
		var ratios []uint
		var accRatio uint
		for i := 0; i < int(n); i++ {
			val := uint(rand.Int63n(amount) + 1)
			if accRatio+val > positive {
				break
			}
			accRatio += val
			ratios = append(ratios, val)
		}

		m := money.New(amount, 1)
		res := m.Allocate(ratios...)
		if sum := money.Sum(res...); amount != sum {
			t.Errorf("split %d by %d, expected %d, got %d", amount, n, amount, sum)
		}
	})
}
