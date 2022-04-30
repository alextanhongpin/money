package money_test

import (
	"math/rand"
	"testing"

	"github.com/alextanhongpin/money"
)

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
