package converter

import (
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GetRandomString returns a random string with length given.
// length < 1 lead to an empty string
// Requires go 1.20 or later -> see https://tip.golang.org/doc/go1.20
// The math/rand package now automatically seeds the global random number generator.
func GetRandomString(length int) string {
	if length < 1 {
		return ""
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(62)]
	}
	return string(b)
}
