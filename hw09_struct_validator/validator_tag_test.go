package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorTag(t *testing.T) {
	t.Run("make tag positive", func(t *testing.T) {
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

	t.Run("tags positive", func(t *testing.T) {
		vts := getValidatorTags("len:11|min:5|max:10")
		require.Len(t, vts, 3)
	})

	t.Run("tags no separator", func(t *testing.T) {
		vts := getValidatorTags("len:11")
		require.Len(t, vts, 1)
	})

	t.Run("tags empty", func(t *testing.T) {
		vts := getValidatorTags("")
		require.Len(t, vts, 0)
	})
}
