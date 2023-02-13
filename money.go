package money

import (
	"errors"
	"fmt"
	"math/big"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

var (
	ErrUnitInvalid      = errors.New("money: unit must be at least 1")
	ErrNegativeAmount   = errors.New("money: negative amount")
	ErrFractionalAmount = errors.New("money: amount cannot be fraction")
)

type Money[T constraints.Integer] struct {
	amount T
	unit   T
}

func (m *Money[T]) Validate() error {
	if m.unit < 1 {
		return fmt.Errorf("%w: %d", ErrUnitInvalid, m.unit)
	}

	if m.amount < 0 {
		return ErrNegativeAmount
	}

	if m.amount%m.unit != 0 {
		return fmt.Errorf("%w: dividing %d by %d", ErrFractionalAmount, m.amount, m.unit)
	}

	return nil
}

func NewMoney[T constraints.Integer](amount, unit T) *Money[T] {
	return &Money[T]{
		amount: amount,
		unit:   unit,
	}
}

func (m *Money[T]) Amount() T {
	return m.amount
}

func (m *Money[T]) Unit() T {
	return m.unit
}

func (m *Money[T]) Split(n uint) []T {
	if err := m.Validate(); err != nil {
		panic(err)
	}

	if n == 0 {
		return make([]T, 0)
	}

	amt := m.amount / m.unit / T(n) * m.unit
	if amt < 0 {
		panic("money: negative amount")
	}

	res := make([]T, n)
	for i := 0; i < int(n)-1; i++ {
		res[i] = amt
	}

	last := m.amount - amt*T(n-1)
	if last < 0 {
		panic("money: negative amount")
	}

	res[len(res)-1] = last

	if m.amount != Sum(res) {
		panic("money: invalid split")
	}

	return res
}

func (m *Money[T]) Allocate(ratios []T) []T {
	if err := m.Validate(); err != nil {
		panic(err)
	}

	if len(ratios) == 0 {
		return make([]T, 0)
	}

	units := m.amount / m.unit

	totalRatios := Sum(ratios)

	n := len(ratios)
	res := make([]T, n)

	total := m.amount
	for i := 0; i < n-1; i++ {
		if ratios[i] == 0 {
			continue
		}

		ratio := divBigRatUint64(uint64(ratios[i]), uint64(totalRatios))
		ratio.Mul(ratio, bigRatFromUint64(uint64(units)))

		i64 := floor(ratio)
		i64.Mul(i64, bigIntFromUint64(uint64(m.unit)))
		u64 := i64.Uint64()

		res[i] = T(u64)
		total -= T(u64)
	}

	res[n-1] = total

	return res
}

func (m *Money[T]) Discount(percent Percent) T {
	if err := m.Validate(); err != nil {
		panic(err)
	}

	if err := percent.Validate(); err != nil {
		panic(err)
	}

	units := m.amount / m.unit

	r64 := divBigRatUint64(uint64(percent), 100)
	r64.Mul(r64, bigRatFromUint64(uint64(units)))

	i64 := ceil(r64)
	i64.Mul(i64, bigIntFromUint64(uint64(m.unit)))

	return T(i64.Uint64())
}

func AllocateMap[T constraints.Ordered, V constraints.Integer](m *Money[V], ratioByKey map[T]V) map[T]V {
	keys := make([]T, 0, len(ratioByKey))
	for k := range ratioByKey {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	ratios := make([]V, len(ratioByKey))
	for i, k := range keys {
		ratios[i] = ratioByKey[k]
	}

	allocations := m.Allocate(ratios)
	res := make(map[T]V)
	for i, k := range keys {
		res[k] = allocations[i]
	}

	return res
}

func divBigRat(a, b *big.Rat) *big.Rat {
	return new(big.Rat).Mul(a, new(big.Rat).Inv(b))
}

func mulBigInt(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

func divBigRatUint64(a, b uint64) *big.Rat {
	return new(big.Rat).SetFrac(
		bigIntFromUint64(a),
		bigIntFromUint64(b),
	)
}

func bigRatFromBigInt(a, b *big.Int) *big.Rat {
	return new(big.Rat).SetFrac(a, b)
}

func bigIntFromUint64(n uint64) *big.Int {
	return new(big.Int).SetUint64(n)
}

func bigRatFromUint64(n uint64) *big.Rat {
	return new(big.Rat).SetUint64(n)
}

// https://golang-nuts.narkive.com/RQ5Nof2y/big-rat-ceil
func ceil(x *big.Rat) *big.Int {
	z := new(big.Int)
	z.Add(x.Num(), x.Denom())
	z.Sub(z, big.NewInt(1))
	z.Div(z, x.Denom())
	return z
}

// Returns a new big.Int set to the floor of x.
func floor(x *big.Rat) *big.Int {
	z := new(big.Int)
	z.Div(x.Num(), x.Denom())
	return z
}
