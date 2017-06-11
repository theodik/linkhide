package shortener

import (
	"bytes"
	"hash/fnv"
)

const (
	offsetBasis32 = 2166136261
	fnvPrime32    = 16777619
	base          = 62
)

var alphabet = [62]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

// Reverse returns its argument string reversed rune-wise left to right.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func fnv1a(str string) int {
	hash := fnv.New32a()
	hash.Write([]byte(str))
	return int(hash.Sum32())
}

func base62(n int) string {
	if n == 0 {
		return alphabet[0]
	}

	var buffer bytes.Buffer

	for n > 0 {
		buffer.WriteString(alphabet[n%base])
		n = n / base
	}

	return reverse(buffer.String())
}

// Encode string into short string
func Encode(str string) string {
	id := fnv1a(str)
	return base62(id)
}
