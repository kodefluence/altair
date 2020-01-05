package util

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func IntToPointer(i int) *int {
	return &i
}

func TimeToPointer(t time.Time) *time.Time {
	return &t
}

func StringToPointer(s string) *string {
	return &s
}

func SHA1() string {
	hasher := sha1.New()
	hasher.Write([]byte(uuid.New().String()))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
