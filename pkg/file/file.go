package file

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveFileHeader(file *multipart.FileHeader, dst string) error {
	// TODO: 去重
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
