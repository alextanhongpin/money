package money

import (
	"math/big"

	"golang.org/x/exp/constraints"
)

func Sum[T constraints.Integer](ns []T) T {
	var tot T
	for i := 0; i < len(ns); i++ {
		tot += ns[i]
	}

	return tot
}

func SumBig(ns []*big.Int) *big.Int {
	tot := big.NewInt(0)
	for i := 0; i < len(ns); i++ {
		tot.Add(tot, ns[i])
	}

	return tot
}
