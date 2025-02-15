package generate

import (
	"crypto/rand"
	"math/big"
)

func Generate(length int) (string, error) {
	var (
		set      string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
		shortUrl []byte = make([]byte, length)
	)
	for i := 0; i < length; i++ {
		r, err := rand.Int(rand.Reader, big.NewInt(63))
		if err != nil {
			return "", err
		}
		shortUrl[i] = set[r.Int64()]
	}
	return string(shortUrl), nil
}
