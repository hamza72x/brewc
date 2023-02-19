package util

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

const randomCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StrContains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)

	for i := range b {
		b[i] = randomCharSet[rand.Intn(len(randomCharSet))]
	}

	return string(b)
}
