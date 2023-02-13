package main

import (
	"fmt"

	"github.com/alextanhongpin/money"
)

func main() {
	usd := USD(50, 30)
	fmt.Println(usd.Split(3))                   // [1676 1676 1678]
	fmt.Println(usd.Allocate([]int64{1, 2, 5})) // [628 1257 3145]
	fmt.Println(usd.Discount(5))                // 252

	sgd := SGD(50, 30)
	fmt.Println(sgd.Split(3))                   // [1675 1675 1680]
	fmt.Println(sgd.Allocate([]int64{1, 2, 5})) // [625 1255 3150]
	fmt.Println(sgd.Discount(5))                // 255

	idr := IDR(534_000)
	fmt.Println(idr.Split(3))                   // [178000 178000 178000]
	fmt.Println(idr.Allocate([]int64{1, 2, 5})) // [66700 133500 333800]
	fmt.Println(idr.Discount(5))                // 26700
}

func USD(dollar, cents int64) *money.Money[int64] {
	cents += dollar * 100
	return money.NewMoney(cents, 1) // 1 penny is the smallest unit.
}

func SGD(dollar, cents int64) *money.Money[int64] {
	// Note that the 5 cents rounding is only valid for offline payment where
	// coins are involved.
	// For digital payments, 1 cent is acceptable.
	cents += dollar * 100
	return money.NewMoney(cents, 5) // 5 cents is the smallest unit.
}

// There are no decimals in Indonesian Rupiah.
func IDR(rupiah int64) *money.Money[int64] {
	return money.NewMoney(rupiah, 100) // 100 rupiah is the smallest unit.
}
