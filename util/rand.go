package util

import "math/rand/v2"

var LettersLower = []rune("abcdefghijklmnopqrstuvwxyz")

func RandomString(n uint64, charPool []rune) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = charPool[rand.IntN(len(charPool))]
	}
	return string(s)
}
