package kaiko

import (
	"math/rand"
	"time"
)

// saltGenerate parameters:
const (
	alpha_bytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	byte_bits     = 6                // 6 bits to represent a letter index
	byte_mask     = 1<<byte_bits - 1 // All 1-bits, as many as byte_bits
	max_alpha_bit = 63 / byte_bits   // # of letter indices fitting in 63 bits
)

// unique rand seed for all salt generation >> based on nano timestamp (cross-compatible??)
var seed = rand.NewSource(time.Now().UnixNano())

// not goroutine safe, needs to instanciate an rand.Int63() per task for concurrent access
func saltGenerate(n int) []byte {
	salt := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for max_alpha_bit characters!
	for i, cache, remain := n-1, seed.Int63(), max_alpha_bit; i >= 0; {
		if remain == 0 {
			cache, remain = seed.Int63(), max_alpha_bit
		}
		if j := int(cache & byte_mask); j < len(alpha_bytes) {
			salt[i] = alpha_bytes[j]
			i--
		}
		cache >>= byte_bits
		remain--
	}
	return salt
}
