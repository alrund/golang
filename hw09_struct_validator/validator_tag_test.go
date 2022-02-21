package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorTag(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		vt := ValidatorTag("len:19")
		name, _ := vt.getName()
		parameter, _ := vt.getParameter()
		require.Equal(t, "len", name)
		require.Equal(t, "19", parameter)
	})

	t.Run("no separator", func(t *testing.T) {
		vt := ValidatorTag("len19")
		_, err := vt.getName()
		require.ErrorIs(t, err, ErrValidatorTagNameSeparatorMissing)
	})

	t.Run("empty name", func(t *testing.T) {
		vt := ValidatorTag(":19")
		_, err := vt.getName()
		require.ErrorIs(t, err, ErrValidatorTagNameEmpty)
	})

	t.Run("empty parameter", func(t *testing.T) {
		vt := ValidatorTag("len:")
		_, err := vt.getParameter()
		require.ErrorIs(t, err, ErrValidatorTagParameterEmpty)
	})
}
