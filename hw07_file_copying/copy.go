package main

import (
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/seth2810/otus_homework/hw07_file_copying/progress"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrLimitNegative         = errors.New("limit should not be negative")
	ErrOffsetNegative        = errors.New("offset should not be negative")
	ErrFromPathIsRequired    = errors.New("source path should not be epmty")
	ErrToPathIsRequired      = errors.New("destination path should not be epmty")
	ErrFromToPathsSame       = errors.New("to path should not be same as from path")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	if fromPath == "" {
		return ErrFromPathIsRequired
	}

	if toPath == "" {
		return ErrToPathIsRequired
	}

	if fromPath == toPath {
		return ErrFromToPathsSame
	}

	if limit < 0 {
		return ErrLimitNegative
	}

	if offset < 0 {
		return ErrOffsetNegative
	}

	var src, dst *os.File

	if src, err = os.Open(fromPath); err != nil {
		return err
	}

	defer src.Close()

	var stat fs.FileInfo

	if stat, err = src.Stat(); err != nil {
		return err
	}

	if stat.Mode().IsDir() {
		return ErrUnsupportedFile
	}

	fileSize := stat.Size()

	if fileSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if dst, err = os.Create(toPath); err != nil {
		return err
	}

	defer dst.Close()

	remainingSize := fileSize - offset

	if limit == 0 || limit > remainingSize {
		limit = remainingSize
	}

	reader := io.NewSectionReader(src, offset, limit)

	bar := progress.NewBar(reader, limit)

	bar.Watch()

	defer bar.Close()

	if _, err = io.CopyN(dst, bar, limit); err != nil {
		return err
	}

	return nil
}
