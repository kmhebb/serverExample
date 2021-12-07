package random

import "math/rand"

func GenerateCode(n int) string {
	chars := "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
	var code []byte
	for i := 0; i < n; i++ {
		code = append(code, chars[rand.Intn(len(chars))])
	}
	return string(code)
}
