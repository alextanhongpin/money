package money

import (
	"errors"
	"fmt"
)

var (
	ErrPercentOutOfRange = errors.New("money: percent must be between 0 and 100")
)

// Percent is a value between 0 and 100.
type Percent uint

// Valid returns true if the range is valid.
func (p Percent) Valid() bool {
	return p.Validate() == nil
}

// Validate checks if the percent range is valid.
func (p *Percent) Validate() error {
	pp := *p
	if pp < 0 || pp > 100 {
		return fmt.Errorf("%w: %d", ErrPercentOutOfRange, pp)
	}

	return nil
}
