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

func PointerToInt(i *int) int {
	if i == nil {
		return 0
	}

	return *i
}

func TimeToPointer(t time.Time) *time.Time {
	return &t
}

func PointerToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return *t
}

func StringToPointer(s string) *string {
	return &s
}

func PointerToString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func SHA1() string {
	hasher := sha1.New()
	hasher.Write([]byte(uuid.New().String()))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
