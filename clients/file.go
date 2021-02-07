package clients

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

type FileClient struct {
	file *os.File
}

func (f *FileClient) WriteString(message string) error {
	length, err := fmt.Fprint(f.file, message)
	if err != nil {
		return fmt.Errorf(
			"could not write to file (%s): %w",
			f.file.Name(),
			err,
		)
	}

	if length < len(message) {
		return fmt.Errorf(
			"could not full message to file (%s)",
			f.file.Name(),
		)
	}

	return nil
}

func (f *FileClient) Close() error {
	err := f.file.Close()
	if err != nil {
		return fmt.Errorf(
			"cannot close file (%s): %w",
			f.file.Name(),
			err,
		)
	}

	return nil
}

func NewFile(uri *url.URL) (*FileClient, error) {
	file, err := os.OpenFile(
		filepath.FromSlash(uri.Path),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not start to (%s): %w",
			uri.String(),
			err,
		)
	}
	return &FileClient{
		file: file,
	}, nil
}
