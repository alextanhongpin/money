package money_test

import (
	"fmt"
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
		unit := uint(rand.Intn(len(fmt.Sprint(amount))) + 1)

		m := money.New(amount, unit)
		msplit := m.Split(uint(n))
		msum := money.Sum(msplit...)

		b := money.NewBig(big.NewInt(amount), unit)
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

		unit := uint(rand.Intn(len(fmt.Sprint(amount))) + 1)

		m := money.New(amount, unit)
		msplit := m.Allocate(ratios...)
		msum := money.Sum(msplit...)

		b := money.NewBig(big.NewInt(amount), unit)
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

func FuzzBigMoneyDiscount(f *testing.F) {
	f.Fuzz(func(t *testing.T, positive uint64) {
		amount := int64(positive)
		percent := money.Percent(rand.Intn(100))

		unit := uint(rand.Intn(len(fmt.Sprint(amount))) + 1)
		m := money.New(amount, unit)
		mper := m.Discount(percent)

		b := money.NewBig(big.NewInt(amount), unit)
		pper := b.Discount(percent)
		if pper.Int64() != mper {
			t.Errorf("discount does not match, %d != %d", pper.Int64(), mper)
		}
	})
}
