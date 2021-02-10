package clients

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"os"
)

type File struct {
	file *os.File
}

func (f *File) ReadString() (string, error) {
	return "", nil
}

func (f *File) WriteString(message string) error {
	length, err := f.file.WriteString(message)
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

func (f *File) Close() error {
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

func NewFile(uri *url.URL) (*File, error) {
	file, err := os.OpenFile(
		uri.Path,
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
	return &File{
		file: file,
	}, nil
}
