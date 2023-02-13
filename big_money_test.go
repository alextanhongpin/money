package money_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/alextanhongpin/money"
	"github.com/stretchr/testify/assert"
)

func TestBigMoney(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		assert := assert.New(t)
		m := money.NewBigMoney(big.NewInt(-1), big.NewInt(1))
		assert.Equal(money.ErrNegativeAmount, m.Validate())
	})

	t.Run("zero", func(t *testing.T) {
		assert := assert.New(t)
		m := money.NewBigMoney(big.NewInt(0), big.NewInt(1))

		assert.Equal(int64(0), m.Amount().Int64())
		assert.Equal(int64(1), m.Unit().Int64())
		for _, amt := range m.Split(3) {
			assert.Equal(int64(0), amt.Int64())
		}

		for _, amt := range m.Allocate([]uint64{1, 2, 3}) {
			assert.Equal(int64(0), amt.Int64())
		}
	})

	t.Run("negative unit", func(t *testing.T) {
		assert := assert.New(t)
		m := money.NewBigMoney(big.NewInt(1), big.NewInt(-1))
		assert.Equal(money.ErrUnitInvalid, m.Validate())
	})

	t.Run("fractional unit", func(t *testing.T) {
		assert := assert.New(t)
		m := money.NewBigMoney(big.NewInt(3), big.NewInt(2))
		assert.True(errors.Is(m.Validate(), money.ErrFractionalAmount))
	})
}

func TestBigMoneySplit(t *testing.T) {
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
		{amount: 100, unit: 5, split: 3, expected: []int64{30, 30, 40}, scenario: "split by 3, unit 5"},
		{amount: 100, unit: 10, split: 3, expected: []int64{30, 30, 40}, scenario: "split by 3, unit 10"},
		{amount: 100, unit: 20, split: 3, expected: []int64{20, 20, 60}, scenario: "split by 3, unit 20"},
		{amount: 100, unit: 50, split: 3, expected: []int64{0, 0, 100}, scenario: "split by 3, unit 50"},
		{amount: 100, unit: 100, split: 3, expected: []int64{0, 0, 100}, scenario: "split by 3, unit 100"},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			assert := assert.New(t)

			m := money.NewBigMoney(big.NewInt(test.amount), big.NewInt(test.unit))
			res := m.Split(test.split)
			for i, exp := range test.expected {
				assert.Equal(big.NewInt(exp), res[i])
			}

			// The amount will only tally if the split is made.
			if test.split > 0 {
				assert.Equal(money.SumBig(res).Uint64(), m.Amount().Uint64())
			}
		})
	}
}

func TestBigMoneyAllocate(t *testing.T) {
	tests := []struct {
		amount   int64
		unit     int64
		allocate []uint64
		expected []int64
		scenario string
	}{
		{amount: 100, unit: 1, allocate: []uint64{}, expected: []int64{}, scenario: "allocate by 0 equally, unit 1"},
		{amount: 100, unit: 1, allocate: []uint64{1, 1, 1}, expected: []int64{33, 33, 34}, scenario: "allocate by 3 equally, unit 1"},
		{amount: 100, unit: 5, allocate: []uint64{1, 1, 1}, expected: []int64{30, 30, 40}, scenario: "allocate by 3 equally, unit 5"},
		{amount: 100, unit: 1, allocate: []uint64{1, 2, 3}, expected: []int64{16, 33, 51}, scenario: "allocate by 3 in ratio 1:2:3, unit 1"},
		{amount: 100, unit: 5, allocate: []uint64{1, 2, 3}, expected: []int64{15, 30, 55}, scenario: "allocate by 3 in ratio 1:2:3, unit 5"},
		{amount: 100, unit: 1, allocate: []uint64{1, 0, 0}, expected: []int64{100, 0, 0}, scenario: "allocate by 3 in ratio 1:0:0, unit 1"},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			assert := assert.New(t)

			m := money.NewBigMoney(big.NewInt(test.amount), big.NewInt(test.unit))
			res := m.Allocate(test.allocate)
			for i, exp := range test.expected {
				assert.Equal(exp, res[i].Int64())
			}

			assert.Equal(big.NewInt(test.amount), m.Amount())
			assert.Equal(big.NewInt(test.unit), m.Unit())
		})
	}
}

func TestBigMoneyDiscount(t *testing.T) {
	tests := []struct {
		amount   int64
		unit     int64
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
			m := money.NewBigMoney(big.NewInt(test.amount), big.NewInt(test.unit))
			assert.Equal(test.expected, m.Discount(money.Percent(test.discount)).Uint64())
		})
	}
}
