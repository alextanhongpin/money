package money

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

type Money struct {
	value int64
	unit  int64
}

func Sum[T constraints.Ordered](ts ...T) T {
	var res T
	for _, t := range ts {
		res += t
	}
	return res
}

func New(value, unit int64) *Money {
	if value <= 0 || unit <= 0 {
		panic("value and unit must be greater than 0")
	}
	return &Money{value, unit}
}

func (m Money) Split(n int) []int64 {
	if n <= 0 {
		panic("cannot split by 0")
	}
	if m.value < int64(n) {
		panic(fmt.Errorf("cannot split %d by %d", m.value, n))
	}

	res := make([]int64, n)
	r := float64(m.value) / float64(m.unit)
	r /= float64(n)

	// TODO: Add other rounding options: Ceil, Floor, Round, Bankers?
	var acc int64
	val := int64(math.Round(r))
	for i := 0; i < n-1; i++ {
		res[i] = val * m.unit
		acc += res[i]
	}

	res[n-1] = m.value - acc
	return res
}

func (m Money) Allocate(ratios ...int64) []int64 {
	var totalRatio int64
	for _, ratio := range ratios {
		totalRatio += ratio
	}

	res := make([]int64, len(ratios))

	var acc int64
	for i := 0; i < len(ratios)-1; i++ {
		r := float64(ratios[i]) / float64(totalRatio)
		r = r * float64(m.value) / float64(m.unit)

		val := int64(math.Round(r))
		res[i] = val * m.unit
		acc += res[i]
	}
	res[len(ratios)-1] = m.value - acc
	return res
}

func (m Money) AllocateMap(ratios map[int64]int64) map[int64]int64 {
	keys := make([]int64, 0, len(ratios))
	for k := range ratios {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	values := make([]int64, len(keys))
	for i, k := range keys {
		values[i] = ratios[k]
	}

	allocations := m.Allocate(values...)

	res := make(map[int64]int64)
	for i, k := range keys {
		res[k] = allocations[i]
	}

	return res
}
