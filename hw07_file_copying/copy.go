package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	ErrThisIsNotAReadableFile = errors.New("this is not a readable file")
	ErrThisIsNotARegularFile  = errors.New("this is not a regular file")
	ErrCreatedFile            = errors.New("unable to create file")
	ErrSeekFile               = errors.New("unable to set the offset")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	size, err := FileSize(fromPath)
	if err != nil {
		fmt.Println(err)
		return ErrUnsupportedFile
	}
	if offset > size {
		return ErrOffsetExceedsFileSize
	}
	if limit > size || limit == 0 {
		limit = size
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return ErrThisIsNotAReadableFile
	}

	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return ErrSeekFile
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return ErrCreatedFile
	}

	err = CopyWithBar(toFile, fromFile, limit)

	fromFile.Close()
	toFile.Close()

	return err
}

func CopyWithBar(toFile, fromFile *os.File, limit int64) error {
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(io.LimitReader(fromFile, limit))
	_, err := io.Copy(toFile, barReader)
	if err != nil {
		return err
	}
	bar.Finish()

	return nil
}

func FileSize(path string) (size int64, err error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	mode := file.Mode()
	if !mode.IsRegular() {
		return 0, ErrThisIsNotARegularFile
	}

	return file.Size(), err
}
