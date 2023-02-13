package money_test

import (
	"errors"
	"testing"

	"github.com/alextanhongpin/money"
	"github.com/stretchr/testify/assert"
)

func TestPercent(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		assert := assert.New(t)

		p := money.Percent(0)
		assert.True(p.Valid())
		assert.Nil(p.Validate())
	})

	t.Run("=100", func(t *testing.T) {
		assert := assert.New(t)

		p := money.Percent(100)
		assert.True(p.Valid())
		assert.Nil(p.Validate())
	})

	t.Run(">100", func(t *testing.T) {
		assert := assert.New(t)

		p := money.Percent(101)
		assert.False(p.Valid())
		assert.True(errors.Is(p.Validate(), money.ErrPercentOutOfRange))
	})
}
