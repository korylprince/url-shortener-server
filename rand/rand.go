package rand

import (
	"log"
	"math/big"
	"math/rand"

	crand "crypto/rand"
)

//Characters is the character-space for random strings
const Characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randInt(max int) int {
	i, err := crand.Int(crand.Reader, big.NewInt(int64(max)))
	if err == nil {
		return int(i.Int64())
	}

	log.Println("WARNING: error using crypto/rand:", err)

	return rand.Intn(max)
}

//String returns a random string of given length composed of characters from Characters
func String(length int) string {
	var s []byte

	for i := 0; i < length; i++ {
		s = append(s, Characters[randInt(len(Characters))])
	}

	return string(s)
}
