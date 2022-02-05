package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const ToFile = "/tmp/xxx"

func TestCopy(t *testing.T) {
	t.Run("get file size", func(t *testing.T) {
		from := "testdata/input.txt"
		var expectedSize int64 = 6617
		fromSize, err := FileSize(from)

		require.Nil(t, err)
		require.Equal(t, expectedSize, fromSize)
	})

	t.Run("get not exists file size", func(t *testing.T) {
		from := "not_exists_file"
		_, err := FileSize(from)
		require.Error(t, err)
	})

	t.Run("copy file", func(t *testing.T) {
		from := "testdata/input.txt"
		var offset, limit int64
		err := Copy(from, ToFile, offset, limit)
		require.Nil(t, err)

		fromFile, _ := os.Stat(from)
		fromSize := fromFile.Size()

		toFile, _ := os.Stat(ToFile)
		toSize := toFile.Size()

		require.Equal(t, fromSize, toSize)
	})

	t.Run("copy a file with limit exceeding the file size", func(t *testing.T) {
		from := "testdata/out_offset0_limit10000.txt"
		var offset int64
		var limit int64 = 10000
		err := Copy(from, ToFile, offset, limit)
		require.Nil(t, err)

		fromFile, _ := os.Stat(from)
		fromSize := fromFile.Size()

		toFile, _ := os.Stat(ToFile)
		toSize := toFile.Size()

		require.Equal(t, fromSize, toSize)
	})

	t.Run("copy a file with offset exceeding the file size", func(t *testing.T) {
		from := "testdata/out_offset6000_limit1000.txt"
		var offset int64 = 6000
		var limit int64 = 1000
		err := Copy(from, ToFile, offset, limit)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	tests1 := []struct {
		title string
		limit int64
	}{
		{title: "out_offset0_limit10.txt", limit: 10},
		{title: "out_offset0_limit1000.txt", limit: 1000},
	}

	for _, tc := range tests1 {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			from := "testdata/" + tc.title
			var offset int64
			err := Copy(from, ToFile, offset, tc.limit)
			require.Nil(t, err)

			toFile, _ := os.Stat(ToFile)
			toSize := toFile.Size()

			require.Equal(t, tc.limit, toSize)
		})
	}

	tests2 := []struct {
		title  string
		offset int64
		limit  int64
	}{
		{title: "out_offset100_limit1000.txt", offset: 100, limit: 1000},
		{title: "out_offset900_limit1000.txt", offset: 900, limit: 1000},
	}

	for _, tc := range tests2 {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			from := "testdata/" + tc.title
			err := Copy(from, ToFile, tc.offset, tc.limit)
			require.Nil(t, err)

			toFile, _ := os.Stat(ToFile)
			toSize := toFile.Size()

			require.Equal(t, tc.limit-tc.offset, toSize)
		})
	}
}
