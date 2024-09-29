package storage

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"
)

func SaveFileFromHeader(header *multipart.FileHeader, p string) error {
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path.Dir(p), 0766); err != nil {
			return err
		}
	}

	src, err := header.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(p)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return err
	}

	return nil
}

func RemoveFileIfExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		return nil
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}
