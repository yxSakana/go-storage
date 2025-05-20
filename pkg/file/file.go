package file

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
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

func CalculateHash(filePath, algo string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var hashVal hash.Hash
	switch algo {
	case "md5":
		hashVal = md5.New()
	case "sha256":
		hashVal = sha256.New()
	default:
		return "", fmt.Errorf("unsupported hash algorithm")
	}

	if _, err := io.Copy(hashVal, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hashVal.Sum(nil)), nil
}

func VerifyFileHash(filepath string, expectedHash string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, file); err != nil {
		return err
	}
	hashBytes := md5Hash.Sum(nil)
	hashStr := hex.EncodeToString(hashBytes)
	if hashStr != expectedHash {
		return fmt.Errorf("file hash %s is not equal to expected hash %s", hashStr, expectedHash)
	}
	return nil
}
