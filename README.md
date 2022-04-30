# money [![Go Reference](https://pkg.go.dev/badge/github.com/alextanhongpin/money.svg)](https://pkg.go.dev/github.com/alextanhongpin/money)

Money package allows splitting and allocation without hanging pennies.

```go
package main

import (
	"fmt"

	"github.com/alextanhongpin/money"
)

func main() {
	usd := USD(50, 30)
	fmt.Println(usd.Split(3))          // [1677 1677 1676]
	fmt.Println(usd.Allocate(1, 2, 5)) // [629 1258 3143]

	sgd := SGD(50, 30)
	fmt.Println(sgd.Split(3))          // [1675 1675 1680]
	fmt.Println(sgd.Allocate(1, 2, 5)) // [630 1260 3140]

	idr := IDR(532_041)                // SGD 50.30 in Rupiah exchange rate.
	fmt.Println(idr.Split(3))          // [177300 177300 177441]
	fmt.Println(idr.Allocate(1, 2, 5)) // [66500 133000 332541]
}

func USD(dollar, cents int64) *money.Money {
	cents += dollar * 100
	return money.New(cents, 1) // 1 penny is the smallest unit.
}

func SGD(dollar, cents int64) *money.Money {
	cents += dollar * 100
	return money.New(cents, 5) // 5 cents is the smallest unit.
}

// There are no decimals in Indonesian Rupiah.
func IDR(rupiah int64) *money.Money {
	return money.New(rupiah, 100) // 100 rupiah is the smallest unit.
}
```
