package crypto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"math/big"
)

const (
	Sha3     = "sha3_24_rounds"
	maxNonce = math.MaxInt64
)

var prefix = []byte("Vega_SPAM_PoW")

// PoW calculates proof of work given block hash, transaction hash, target difficulty and a hash function.
// returns the nonce, the hash and th error if any.
func PoW(blockHash string, txID string, difficulty uint, hashFunction string) (uint64, []byte, error) {
	var hashInt big.Int
	var h []byte
	var err error
	nonce := uint64(0)

	if difficulty > 256 {
		return nonce, h, errors.New("invalid difficulty")
	}

	target := big.NewInt(1)
	target.Lsh(target, 256-difficulty)

	if len(txID) < 1 {
		return nonce, h, errors.New("transaction ID cannot be empty")
	}

	if len(blockHash) != 64 {
		return nonce, h, errors.New("incorrect block hash")
	}

	for nonce < maxNonce {
		data := prepareData(blockHash, txID, nonce)
		h, err = hash(data, hashFunction)
		if err != nil {
			return nonce, h, err
		}

		hashInt.SetBytes(h[:])
		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, h[:], nil
}

// Verify checks that the hash with the given nonce results in the target difficulty.
func Verify(blockHash string, txID string, nonce uint64, hashFuncion string, difficulty uint) (bool, big.Int) {
	var hashInt big.Int
	target := big.NewInt(1)

	if difficulty > 256 {
		return false, hashInt
	}

	target.Lsh(target, 256-difficulty)

	if len(txID) < 1 {
		return false, hashInt
	}

	if len(blockHash) != 64 {
		return false, hashInt
	}

	data := prepareData(blockHash, txID, nonce)
	h, err := hash(data, hashFuncion)
	if err != nil {
		return false, hashInt
	}
	hashInt.SetBytes(h[:])
	return hashInt.Cmp(target) == -1, hashInt
}

func prepareData(blockHash string, txID string, nonce uint64) []byte {
	data := bytes.Join(
		[][]byte{
			prefix,
			[]byte(blockHash),
			[]byte(txID),
			IntToHex(nonce),
		},
		[]byte{},
	)

	return data
}

func hash(data []byte, hashFunction string) ([]byte, error) {
	if hashFunction == Sha3 {
		return Hash(data), nil
	}
	return nil, errors.New("unknown hash function")
}

func IntToHex(num uint64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
