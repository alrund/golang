package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	TestDataDir = "testdata/env"
)

func TestReadDir(t *testing.T) {
	env, _ := ReadDir(TestDataDir)

	t.Run("Simple", func(t *testing.T) {
		require.Equal(t, "bar", env["BAR"].Value)
		require.Equal(t, false, env["BAR"].NeedRemove)
	})

	t.Run("Quoted", func(t *testing.T) {
		require.Equal(t, "\"hello\"", env["HELLO"].Value)
		require.Equal(t, false, env["HELLO"].NeedRemove)
	})

	t.Run("With line break", func(t *testing.T) {
		require.Equal(t, "   foo\nwith new line", env["FOO"].Value)
		require.Equal(t, false, env["FOO"].NeedRemove)
	})

	t.Run("Empty", func(t *testing.T) {
		require.Equal(t, "", env["EMPTY"].Value)
		require.Equal(t, false, env["EMPTY"].NeedRemove)
	})

	t.Run("Unset", func(t *testing.T) {
		require.Equal(t, "", env["UNSET"].Value)
		require.Equal(t, true, env["UNSET"].NeedRemove)
	})

	t.Run("Name with equal", func(t *testing.T) {
		require.NotContains(t, env, "ZZZ=XXX")
	})
}
