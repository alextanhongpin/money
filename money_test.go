package money_test

import (
	"math/rand"
	"testing"

	"github.com/alextanhongpin/money"
)

func FuzzMoneySplit(f *testing.F) {
	f.Fuzz(func(t *testing.T, amount int64) {
		if amount <= 0 {
			return
		}
		n := rand.Int63n(amount)
		if n <= 0 {
			return
		}
		m := money.New(amount, 1)
		res := m.Split(int(n))
		if sum := money.Sum(res...); amount != sum {

			t.Errorf("split %d by %d, expected %d, got %d", amount, n, amount, sum)
		}
	})
}

func FuzzMoneyAllocate(f *testing.F) {
	f.Fuzz(func(t *testing.T, amount int64) {
		if amount <= 0 {
			return
		}
		n := rand.Int63n(amount)
		if n <= 0 {
			return
		}

		var ratios []int64
		var accRatio int64
		for i := 0; i < int(n); i++ {
			val := rand.Int63n(amount)
			if accRatio+val > amount {
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
