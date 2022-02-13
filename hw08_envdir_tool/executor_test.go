package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Empty cmd", func(t *testing.T) {
		returnCode := RunCmd([]string{}, Environment{})
		require.Equal(t, 1, returnCode)
	})
	t.Run("Arguments", func(t *testing.T) {
		returnCode := RunCmd([]string{"testdata/echo.sh", "arg1", "arg1"}, Environment{})
		require.Equal(t, 0, returnCode)
	})
	t.Run("Exit code", func(t *testing.T) {
		returnCode := RunCmd([]string{"ls", "non-exists"}, Environment{})
		require.Equal(t, 2, returnCode)
	})
}
