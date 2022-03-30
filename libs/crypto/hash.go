package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

func Hash(key []byte) []byte {
	hasher := sha3.New256()
	hasher.Write(key)
	return hasher.Sum(nil)
}

func RandomHash() string {
	data := make([]byte, 10)
	rand.Read(data)
	return strings.ToUpper(hex.EncodeToString(Hash(data)))
}

func AsSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
