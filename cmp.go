package money

import "math/big"

func isLt(a, b *big.Int) bool {
	return -1 == a.Cmp(b)
}

func isEq(a, b *big.Int) bool {
	return 0 == a.Cmp(b)
}
