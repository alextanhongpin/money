package money_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/alextanhongpin/money"
)

func FuzzBigMoneySplit(f *testing.F) {
	f.Fuzz(func(t *testing.T, positive uint64, n uint64) {
		if positive < n {
			positive, n = n, positive
		}
		positive++
		n++

		amount := int64(positive)

		m := money.New(amount, 1)
		msplit := m.Split(uint(n))
		msum := money.Sum(msplit...)

		b := money.NewBig(big.NewInt(amount), 1)
		bsplit := b.Split(uint(n))
		bsum := money.SumBig(bsplit...).Int64()

		if amount != msum {
			t.Errorf("split %d by %d, expected %d, got %d", amount, n, amount, msum)
		}

		if bsum != msum {
			t.Errorf("split %d by %d, expected %d, got %d", amount, n, msum, bsum)
		}
	})
}

func FuzzBigMoneyAllocate(f *testing.F) {
	f.Fuzz(func(t *testing.T, positive, n uint64) {
		if positive < n {
			positive, n = n, positive
		}
		positive++
		n++

		amount := int64(positive)

		var ratios []uint
		var total uint
		for i := 0; i < int(n); i++ {
			val := uint(rand.Int63n(amount) + 1)
			if total+val > uint(amount) {
				break
			}
			total += val
			ratios = append(ratios, val)
		}

		m := money.New(amount, 1)
		msplit := m.Allocate(ratios...)
		msum := money.Sum(msplit...)

		b := money.NewBig(big.NewInt(amount), 1)
		bsplit := b.Allocate(ratios...)
		bsum := money.SumBig(bsplit...).Int64()

		if amount != msum {
			t.Errorf("split %d by %d, expected %d, got %d", amount, n, amount, msum)
		}

		if bsum != msum {
			t.Errorf("split %d by %d, expected %d, got %d", amount, n, msum, bsum)
		}
	})
}
