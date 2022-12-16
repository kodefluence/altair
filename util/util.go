package util

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
)

func IntToPointer(i int) *int {
	return &i
}

func ReadFileContent(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func PointerToInt(i *int) interface{} {
	if i == nil {
		return nil
	}

	return *i
}

func TimeToPointer(t time.Time) *time.Time {
	return &t
}

func PointerToTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}

	return *t
}

func StringToPointer(s string) *string {
	return &s
}

func PointerToString(s *string) interface{} {
	if s == nil {
		return nil
	}

	return *s
}

func SHA1() string {
	hasher := sha1.New()
	hasher.Write([]byte(uuid.New().String()))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
