package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func NewRandomString(aliasLength int) (string, error) {

	var latters = []rune("abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789")
	alias := make([]rune, aliasLength)

	for i := range alias {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(latters))))
		if err != nil {
			return "", fmt.Errorf("failed generate random int: %w", err)
		}
		alias[i] = latters[n.Int64()]
	}

	return string(alias), nil
}
