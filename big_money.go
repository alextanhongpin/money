package money

import (
	"fmt"
	"math/big"
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
)

type BigMoney struct {
	amount *big.Int
	unit   *big.Int
}

func (m *BigMoney) Validate() error {
	if isLt(m.amount, zero) {
		return ErrNegativeAmount
	}

	if isLt(m.unit, one) {
		return ErrUnitInvalid
	}

	mod := new(big.Int)
	_, _ = new(big.Int).DivMod(m.amount, m.unit, mod)
	if !isEq(mod, zero) {
		return fmt.Errorf("%w: dividing %d by %d", ErrFractionalAmount, m.amount, m.unit)
	}

	return nil
}

func NewBigMoney(amount, unit *big.Int) *BigMoney {
	return &BigMoney{
		amount: amount,
		unit:   unit,
	}
}

func (m *BigMoney) Amount() *big.Int {
	return new(big.Int).Set(m.amount)
}

func (m *BigMoney) Unit() *big.Int {
	return new(big.Int).Set(m.unit)
}

func (m *BigMoney) Split(n uint) []*big.Int {
	if err := m.Validate(); err != nil {
		panic(err)
	}

	if n == 0 {
		return make([]*big.Int, 0)
	}

	nRat64 := new(big.Rat).SetInt64(int64(n))
	units := bigRatFromBigInt(m.amount, m.unit)

	per := floor(divBigRat(units, nRat64))
	res := make([]*big.Int, n)

	total := m.Amount()
	for i := 0; i < int(n)-1; i++ {
		i64 := mulBigInt(per, m.unit)
		res[i] = i64
		total.Sub(total, i64)
	}

	res[len(res)-1] = total

	return res
}

func (m *BigMoney) Allocate(ratios []uint64) []*big.Int {
	if err := m.Validate(); err != nil {
		panic(err)
	}

	if len(ratios) == 0 {
		return make([]*big.Int, 0)
	}

	// Find the total ratio first.
	var ratio uint64
	for _, r := range ratios {
		ratio += r
	}

	ratioInt64 := bigRatFromUint64(ratio)
	unitsInt64 := bigRatFromBigInt(m.amount, m.unit)

	n := len(ratios)
	res := make([]*big.Int, n)

	total := m.Amount()
	for i := 0; i < n-1; i++ {
		if ratios[i] == 0 {
			res[i] = big.NewInt(0)
			continue
		}

		// shares = ratio / total_ratio
		r64 := divBigRat(divBigRatUint64(ratios[i], 1), ratioInt64)

		// shares = shares * (amount / unit)
		r64.Mul(r64, unitsInt64)

		// shares = floor(shares)
		i64 := floor(r64)

		// shares = shares * m.unit
		i64 = mulBigInt(i64, m.unit)

		res[i] = i64
		total.Sub(total, i64)
	}

	res[n-1] = total

	return res
}

func (m *BigMoney) Discount(percent Percent) *big.Int {
	if err := m.Validate(); err != nil {
		panic(err)
	}

	if err := percent.Validate(); err != nil {
		panic(err)
	}

	units := bigRatFromBigInt(m.amount, m.unit)

	r64 := divBigRatUint64(uint64(percent), 100)
	r64.Mul(r64, units)

	i64 := ceil(r64)
	i64.Mul(i64, m.unit)

	return i64
}
