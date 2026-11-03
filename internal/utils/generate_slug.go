package utils

import (
	"math/rand"
	"strings"
)

func GenerateSlug() string {
	var slug strings.Builder

	char := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz"
	lenSlug := 6

	for range lenSlug {
		slug.WriteString(string(char[rand.Intn(lenSlug)]))
	}

	return slug.String()
}
