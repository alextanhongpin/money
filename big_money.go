package money

import (
	"fmt"
	"math/big"
	"sort"

	"golang.org/x/exp/constraints"
)

// BigMoney represents a monetary unit in big integers.
type BigMoney struct {
	value *big.Int
	unit  uint
}

// NewBig returns a pointer to BigMoney.
func NewBig(value *big.Int, unit uint) *BigMoney {
	return &BigMoney{
		value: value,
		unit:  unit,
	}
}

// Value returns a copy of the amount.
func (b BigMoney) Value() *big.Int {
	return new(big.Int).Set(b.value)
}

// Split splits the amount between n. Panics if n is zero or if the amount is less than n.
func (b BigMoney) Split(n uint) []*big.Int {
	if n == 0 {
		panic("BigMoneySplitErr: n must be greater than 0")
	}
	if b.value.Cmp(big.NewInt(int64(n))) == -1 {
		panic(fmt.Errorf("BigMoneySplitErr: cannot split %d by %d", b.value, n))
	}

	r := new(big.Rat).SetInt(b.value)
	r = r.Mul(r, big.NewRat(1, int64(n*b.unit)))
	e := bigRatToBigInt(r)
	e = e.Mul(e, big.NewInt(int64(b.unit)))

	res := make([]*big.Int, n)
	for i := 0; i < int(n)-1; i++ {
		res[i] = new(big.Int).Set(e)
	}
	rem := new(big.Int).Set(b.value)
	res[int(n)-1] = rem.Sub(rem, e.Mul(e, big.NewInt(int64(n)-1)))

	if SumBig(res...).Cmp(b.value) != 0 {
		panic("BigMoneySplitErr: allocated amount doesnot sum up to total amount")
	}

	return res
}

// Allocate attempts to distribute the amount based on the ratios given.
func (b BigMoney) Allocate(ratios ...uint) []*big.Int {
	var total uint
	for _, r := range ratios {
		total += r
	}

	res := make([]*big.Int, len(ratios))
	acc := new(big.Int)
	for i := 0; i < len(ratios)-1; i++ {
		ratio := int64(ratios[i])

		r := new(big.Rat).SetInt(b.value)
		r.Mul(r, big.NewRat(ratio, int64(total*b.unit)))

		e := bigRatToBigInt(r)
		e = e.Mul(e, big.NewInt(int64(b.unit)))

		res[i] = e
		acc.Add(acc, e)
	}

	rem := new(big.Int).Set(b.value)
	res[len(ratios)-1] = rem.Sub(rem, acc)

	if SumBig(res...).Cmp(b.value) != 0 {
		panic("BigMoneyAllocateErr: allocated amount doesnot sum up to total amount")
	}

	return res
}

func (b BigMoney) Discount(percent Percent) *big.Int {
	rat := new(big.Rat).SetInt(b.value)
	rat.Mul(rat, big.NewRat(int64(percent), 100*int64(b.unit)))

	res := ratCeil(rat)
	res.Mul(res, big.NewInt(int64(b.unit)))
	return res
}

// AllocateBigMap is a convenient method to perform
// consistent allocations based on ordered keys.
func AllocateBigMap[T constraints.Ordered](m *BigMoney, ratios map[T]uint) map[T]*big.Int {
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

	res := make(map[T]*big.Int)
	for i, k := range keys {
		res[k] = allocations[i]
	}

	return res
}

func bigRatToBigInt(r *big.Rat) *big.Int {
	i := new(big.Int)
	s := r.FloatString(0)
	fmt.Sscan(s, i)
	return i
}

// SumBig is a convenient method to sum values.
func SumBig(ints ...*big.Int) *big.Int {
	acc := new(big.Int)
	for _, n := range ints {
		acc.Add(acc, n)
	}
	return acc
}

// Returns a new big.Int set to the ceiling of x.
func ratCeil(x *big.Rat) *big.Int {
	z := new(big.Int)
	m := new(big.Int)
	z.DivMod(x.Num(), x.Denom(), m)
	if m.Cmp(big.NewInt(0)) == 1 {
		z.Add(z, big.NewInt(1))
	}
	return z
}
