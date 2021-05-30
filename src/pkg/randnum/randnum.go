package randnum

import (
	"crypto/rand"
	"io"
)

func RandNumString(length int) string {
	var nums = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = nums[int(b[i])%len(nums)]
	}
	return string(b)
}
