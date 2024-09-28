package storage

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveFileFromHeader(header *multipart.FileHeader, path string) error {
	src, err := header.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(path)
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
