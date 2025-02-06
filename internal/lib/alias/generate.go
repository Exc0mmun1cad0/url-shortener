package alias

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base62         = 62
)

// Generate creates alias of length aliasLength for given url.
func Generate(input string, aliasLength int) (string, error) {
	timestamp := time.Now().UnixNano() // salt for sha256 hashing source url
	hash := sha256.Sum256(
		[]byte(input + strconv.Itoa(int(timestamp))),
	)
	hashString := fmt.Sprintf("%x", hash)

	base62Hash, err := base62Encode(hashString, 16)
	if err != nil {
		return "", fmt.Errorf("cannot encode url hash into base62: %w", err)
	}

	// Only a part of hash is returned in order to make alias shorter.
	return base62Hash[:aliasLength], nil
}

func base62Encode(input string, fromBase int) (string, error) {
	src := new(big.Int)
	src.SetString(input, fromBase)

	var result strings.Builder
	base62 := big.NewInt(base62)
	zero := big.NewInt(0)
	for src.Cmp(zero) == 1 {
		newDigit := new(big.Int)
		src.DivMod(src, base62, newDigit)

		err := result.WriteByte(base62Alphabet[newDigit.Int64()])
		if err != nil {
			return "", fmt.Errorf("cannot add digit into result number: %w", err)
		}
	}

	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}
