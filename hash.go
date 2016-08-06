package main

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func generateHash(s string) string {
	hasher := sha1.New()
	_, _ = io.WriteString(hasher, s)
	return fmt.Sprintf("%X", string(hasher.Sum(nil)))
}
