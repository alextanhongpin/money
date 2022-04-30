package money

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

// Money represents a monetary unit.
type Money struct {
	value int64
	unit  uint
}

// New returns a pointer to Money.
func New(value int64, unit uint) *Money {
	return &Money{value, unit}
}

// Value returns the amount.
func (m Money) Value() int64 {
	return m.value
}

// Split splits the amount between n. Panics if n is zero or if the amount is less than n.
func (m Money) Split(n uint) []int64 {
	if n == 0 {
		panic("MoneySplitErr: n must be greater than 0")
	}
	if m.value < int64(n) {
		panic(fmt.Errorf("MoneySplitErr: cannot split %d by %d", m.value, n))
	}

	res := make([]int64, n)
	r := float64(m.value) / float64(n*m.unit)

	// TODO: Add other rounding options: Ceil, Floor, Round, Bankers?
	var acc int64
	val := int64(math.Round(r))
	for i := 0; i < int(n)-1; i++ {
		res[i] = val * int64(m.unit)
		acc += res[i]
	}

	res[n-1] = m.value - acc

	if Sum(res...) != m.value {
		panic("MoneySplitErr: split amount does not add up to total amount")
	}

	return res
}

// Allocate attempts to distribute the amount based on the ratios given.
func (m Money) Allocate(ratios ...uint) []int64 {
	if len(ratios) == 0 {
		panic("MoneyAllocateErr: ratios is required")
	}

	var totalRatio uint
	for _, ratio := range ratios {
		totalRatio += ratio
	}

	res := make([]int64, len(ratios))

	var acc int64
	for i := 0; i < len(ratios)-1; i++ {
		r := float64(ratios[i]) / float64(totalRatio)
		r = r * float64(m.value) / float64(m.unit)

		val := int64(math.Round(r))
		res[i] = val * int64(m.unit)
		acc += res[i]
	}
	res[len(ratios)-1] = m.value - acc

	if Sum(res...) != m.value {
		panic("MoneyAllocateErr: allocated amount does not add up to total amount")
	}

	return res
}

// Discount returns the discounted amount that greater or equal the percent
// discount.
func (m Money) Discount(percent Percent) int64 {
	ratio := float64(m.value) * float64(percent) / float64(100*m.unit)
	return int64(math.Ceil(ratio) * float64(m.unit))
}

// AllocateMap is a convenient method to perform
// consistent allocations based on ordered keys.
func AllocateMap[T constraints.Ordered](m *Money, ratios map[T]uint) map[T]int64 {
	keys := make([]T, 0, len(ratios))
	for k := range ratios {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	values := make([]uint, len(keys))
	for i, k := range keys {
		values[i] = ratios[k]
	}

	allocations := m.Allocate(values...)

	res := make(map[T]int64)
	for i, k := range keys {
		res[k] = allocations[i]
	}

	return res
}

// Sum is a convenient method to sum values.
func Sum[T constraints.Ordered](ts ...T) T {
	var res T
	for _, t := range ts {
		res += t
	}
	return res
}
