package money_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/alextanhongpin/money"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func ExampleSplit() {
	m := money.NewMoney(5030, 1)
	s := m.Split(3)

	fmt.Println(s, money.Sum(s))
	// Output: [1676 1676 1678] 5030
}

func ExampleAllocate() {
	m := money.NewMoney(5030, 1)
	a := m.Allocate([]int{1, 2, 5})

	fmt.Println(a, money.Sum(a))
	// Output: [628 1257 3145] 5030
}

func ExampleAllocateMapIntKey() {
	m := money.NewMoney(5030, 1)
	a := money.AllocateMap(m, map[int64]int{
		1000: 1,
		2000: 2,
		3000: 5,
	})

	values := make([]int, 0, len(a))
	for _, val := range a {
		values = append(values, val)
	}

	fmt.Println(a, money.Sum(values))
	// Output: map[1000:628 2000:1257 3000:3145] 5030
}

func ExampleAllocateMapStrKey() {
	m := money.NewMoney(5030, 1)
	a := money.AllocateMap(m, map[string]int{
		"a": 5,
		"b": 2,
		"c": 1,
	})

	values := make([]int, 0, len(a))
	for _, val := range a {
		values = append(values, val)
	}

	fmt.Println(a, money.Sum(values))
	// Output: map[a:3143 b:1257 c:630] 5030
}

func TestMoney(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		assert := assert.New(t)

		m := money.NewMoney(-1, 1)
		assert.Equal(money.ErrNegativeAmount, m.Validate())
	})

	t.Run("zero", func(t *testing.T) {
		assert := assert.New(t)

		m := money.NewMoney(0, 1)
		assert.Equal(0, m.Amount())
		assert.Equal(1, m.Unit())
		assert.Equal([]int{0, 0, 0}, m.Split(3))
		assert.Equal([]int{0, 0, 0}, m.Allocate([]int{1, 2, 3}))
	})

	t.Run("one", func(t *testing.T) {
		assert := assert.New(t)

		m := money.NewMoney(1, 1)
		assert.Equal(1, m.Amount())
		assert.Equal(1, m.Unit())
		assert.Equal([]int{0, 0, 1}, m.Split(3))
		assert.Equal([]int{0, 0, 1}, m.Allocate([]int{1, 2, 3}))
	})

	t.Run("negative unit", func(t *testing.T) {
		assert := assert.New(t)

		m := money.NewMoney(1, -1)
		assert.True(errors.Is(m.Validate(), money.ErrUnitInvalid))
	})

	t.Run("fractional unit", func(t *testing.T) {
		assert := assert.New(t)

		m := money.NewMoney(3, 2)
		assert.True(errors.Is(m.Validate(), money.ErrFractionalAmount))
	})
}

func TestMoneySplit(t *testing.T) {
	tests := []struct {
		amount   int64
		unit     int64
		split    uint
		expected []int64
		scenario string
	}{
		{amount: 100, unit: 1, split: 0, expected: []int64{}, scenario: "split by 0, unit 1"},
		{amount: 100, unit: 1, split: 1, expected: []int64{100}, scenario: "split by 1, unit 1"},
		{amount: 100, unit: 1, split: 2, expected: []int64{50, 50}, scenario: "split by 2, unit 1"},
		{amount: 100, unit: 1, split: 3, expected: []int64{33, 33, 34}, scenario: "split by 3, unit 1"},
		{amount: 100, unit: 1, split: 3, expected: []int64{33, 33, 34}, scenario: "split by 3, unit 1"},
		{amount: 100, unit: 5, split: 3, expected: []int64{30, 30, 40}, scenario: "split by 3, unit 5"},
		{amount: 100, unit: 10, split: 3, expected: []int64{30, 30, 40}, scenario: "split by 3, unit 10"},
		{amount: 100, unit: 20, split: 3, expected: []int64{20, 20, 60}, scenario: "split by 3, unit 20"},
		{amount: 100, unit: 50, split: 3, expected: []int64{0, 0, 100}, scenario: "split by 3, unit 50"},
		{amount: 100, unit: 100, split: 3, expected: []int64{0, 0, 100}, scenario: "split by 3, unit 100"},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			assert := assert.New(t)

			m := money.NewMoney(test.amount, test.unit)
			res := m.Split(test.split)
			assert.Equal(test.expected, res)
			assert.Equal(test.amount, m.Amount())
			assert.Equal(test.unit, m.Unit())
		})
	}
}

func TestMoneyAllocate(t *testing.T) {
	tests := []struct {
		amount   uint64
		unit     uint64
		allocate []uint64
		expected []uint64
		scenario string
	}{
		{amount: 100, unit: 1, allocate: []uint64{1, 1, 1}, expected: []uint64{33, 33, 34}, scenario: "allocate by 3 equally, unit 1"},
		{amount: 100, unit: 5, allocate: []uint64{1, 1, 1}, expected: []uint64{30, 30, 40}, scenario: "allocate by 3 equally, unit 5"},
		{amount: 100, unit: 1, allocate: []uint64{1, 2, 3}, expected: []uint64{16, 33, 51}, scenario: "allocate by 3 in ratio 1:2:3, unit 1"},
		{amount: 100, unit: 5, allocate: []uint64{1, 2, 3}, expected: []uint64{15, 30, 55}, scenario: "allocate by 3 in ratio 1:2:3, unit 5"},
		{amount: 100, unit: 1, allocate: []uint64{1, 0, 0}, expected: []uint64{100, 0, 0}, scenario: "allocate by 3 in ratio 1:0:0, unit 1"},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			assert := assert.New(t)

			m := money.NewMoney(test.amount, test.unit)
			res := m.Allocate(test.allocate)
			assert.Equal(test.expected, res)
		})
	}
}

func TestMoneyAllocateMap(t *testing.T) {
	t.Run("allocateMap by 3 equally, unit 1", func(t *testing.T) {
		assert := assert.New(t)

		m := money.NewMoney(100, 1)
		res := money.AllocateMap(m, map[string]int{
			"a": 1,
			"c": 1,
			"b": 1,
		})

		assert.Equal(33, res["a"])
		assert.Equal(33, res["b"])
		assert.Equal(34, res["c"])
	})
}

func TestMoneyDiscount(t *testing.T) {
	tests := []struct {
		amount   uint64
		unit     uint64
		discount uint
		expected uint64
		scenario string
	}{
		{amount: 100, unit: 1, discount: 5, expected: 5, scenario: "5% discount from 100, unit 1"},
		{amount: 100, unit: 5, discount: 5, expected: 5, scenario: "5% discount from 100, unit 5"},
		{amount: 100, unit: 5, discount: 11, expected: 15, scenario: "5% discount from 100, unit 5"},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			assert := assert.New(t)
			m := money.NewMoney(test.amount, test.unit)
			assert.Equal(test.expected, m.Discount(money.Percent(test.discount)))
		})
	}
}

func FuzzSplit(f *testing.F) {
	f.Fuzz(func(t *testing.T, amount uint64, n uint) {
		// Ensure the amount is greater than split.
		if amount < uint64(n) {
			amount, n = uint64(n), uint(amount)
		}

		// Ensure the split is at least one, so that the
		// amount tally.
		if n == 0 {
			n = 1
		}

		{
			m := money.NewMoney(amount, 1)
			res := m.Split(n)

			// Ensure that the split amount sums up to the original amount.
			if exp, got := m.Amount(), money.Sum(res); exp != got {
				t.Errorf("split %d by %d, exp %d, got %d", amount, n, exp, got)
			}
		}

		{
			b := money.NewBigMoney(new(big.Int).SetUint64(amount), big.NewInt(1))
			res := b.Split(n)

			// Ensure that the split amount sums up to the original amount.
			if exp, got := b.Amount().Uint64(), money.SumBig(res).Uint64(); exp != got {
				t.Errorf("split %d by %d, exp %d, got %d", amount, n, exp, got)
			}
		}
	})
}

func FuzzAllocate(f *testing.F) {
	f.Fuzz(func(t *testing.T, amount uint64) {
		// Ensure that the total ratio does not exceed the amount.
		nratio := rand.Intn(10)
		ratios := make([]uint64, nratio)
		vratio := amount
		for i := 0; i < nratio; i++ {
			n := rand.Uint64()
			if vratio-n > 0 {
				ratios[i] = n
				vratio -= n
				continue
			}
			break
		}

		// Ensure the allocation is done so that the amount tally.
		if len(ratios) == 0 {
			ratios = make([]uint64, 1)
			ratios[0] = 1
		}

		{
			m := money.NewMoney(amount, 1)
			res := m.Allocate(ratios)

			// Ensure that the split amount sums up to the original amount.
			if exp, got := m.Amount(), money.Sum(res); exp != got {
				t.Errorf("allocate %d by %v, exp %d, got %d", amount, ratios, exp, got)
			}
		}

		{
			m := money.NewBigMoney(new(big.Int).SetUint64(amount), big.NewInt(1))
			res := m.Allocate(ratios)

			// Ensure that the split amount sums up to the original amount.
			if exp, got := m.Amount().Uint64(), money.SumBig(res).Uint64(); exp != got {
				t.Errorf("allocate big %d by %v, exp %d, got %d", amount, ratios, exp, got)
			}
		}
	})
}

func FuzzDiscount(f *testing.F) {
	f.Fuzz(func(t *testing.T, amount uint64) {
		discount := money.Percent(rand.Intn(101))

		m := money.NewMoney(amount, 1)
		b := money.NewBigMoney(new(big.Int).SetUint64(amount), big.NewInt(1))

		// Ensure that the split amount sums up to the original amount.
		if exp, got := m.Discount(discount), b.Discount(discount).Uint64(); exp != got {
			t.Errorf("%d discount for %d, exp %d, got %d", discount, amount, exp, got)
		}
	})
}
