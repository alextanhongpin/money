# money [![Go Reference](https://pkg.go.dev/badge/github.com/alextanhongpin/money.svg)](https://pkg.go.dev/github.com/alextanhongpin/money)

Money package allows splitting and allocation without hanging pennies.

```go
func main() {
	cents := int64(5030) // 5030 cents, 50.30 USD
	unit := int64(1)     // 1 cent is the smallest amount divisible.

	fmt.Println("Split 50.30 USD by 3")
	m := money.New(cents, unit)
	split := m.Split(3)
	fmt.Println(split, money.Sum(split...))

	fmt.Println("\nAlloc 50.30 USD between 1, 2 and 5")
	alloc := m.Allocate(1, 2, 5)
	fmt.Println(alloc, money.Sum(alloc...))

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
```
