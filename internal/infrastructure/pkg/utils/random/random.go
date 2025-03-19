package random

import (
	"math/rand"
	"time"
)

func GenerateSixDigitOtp() int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := 100000 + rng.Intn(900000)
	return otp
}

func GenerateVoucherCode(prefix string) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	codeLength := 10
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, codeLength)

	for i := range code {
		code[i] = characters[r.Intn(len(characters))]
	}

	return prefix + string(code)
}
