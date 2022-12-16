package util

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
)

type Value interface {
	int | string | time.Time
}

func ValueToPointer[V Value](v V) *V {
	return &v
}

func PointerToValue[V Value](v *V) V {
	if v == nil {
		return *new(V)
	}
	return *v
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

func SHA1() string {
	hasher := sha1.New()
	hasher.Write([]byte(uuid.New().String()))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
