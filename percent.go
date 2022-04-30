package money

// Percent is a value between 0 and 100.
type Percent uint

func (p Percent) Valid() bool {
	return p >= 0 && p <= 100
}
func NewPercent(value uint) Percent {
	p := Percent(value)
	if !p.Valid() {
		panic("percent must be between 0 and 100")
	}
	return p
}
