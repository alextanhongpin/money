package main

import (
	"fmt"

	"github.com/alextanhongpin/money"
)

func main() {
	cents := int64(5030) // 5030 cents, 50.30 USD
	unit := uint(1)      // 1 cent is the smallest amount divisible.

	fmt.Println("Split 50.30 USD by 3")
	m := money.New(cents, unit)
	split := m.Split(3)
	fmt.Println(split, money.Sum(split...))

	fmt.Println("\nAlloc 50.30 USD between 1, 2 and 5")
	alloc := m.Allocate(1, 2, 5)
	fmt.Println(alloc, money.Sum(alloc...))

	{
		allocMapInt64 := money.AllocateMap(m, map[int64]uint{
			1000: 1,
			2000: 2,
			3000: 5,
		})
		values := make([]int64, 0, len(allocMapInt64))
		for _, val := range allocMapInt64 {
			values = append(values, val)
		}
		fmt.Println(allocMapInt64, money.Sum(values...))
	}

	{
		allocMapStr := money.AllocateMap(m, map[string]uint{
			"a": 5,
			"b": 2,
			"c": 1,
		})
		values := make([]int64, 0, len(allocMapStr))
		for _, val := range allocMapStr {
			values = append(values, val)
		}
		fmt.Println(allocMapStr, money.Sum(values...))
	}

	{
		fmt.Println("\nSplit 50.30 SGD by 3")
		// 50.30 SGD
		sgd := SGD(5030)
		split := sgd.Split(3)
		fmt.Println(split, money.Sum(split...))

		fmt.Println("\nAlloc 50.30 SGD between 1, 2 and 5")
		alloc := sgd.Allocate(1, 2, 5)
		fmt.Println(alloc, money.Sum(alloc...))
	}
}

func SGD(cents int64) *money.Money {
	// Smallest coin is 5 cents.
	return money.New(cents, 5)
}
