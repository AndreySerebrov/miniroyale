package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

const dataBlockLen = 36
const targetBits = 24

func GetTargetButs() int {
	return targetBits
}

func GetHash(exit chan struct{}, data []byte) (big.Int, uint64, error) {
	nonce := uint64(0)
	var hashInt big.Int
	var ok bool
	var err error
	if len(data) != dataBlockLen {
		return hashInt, 0, fmt.Errorf("expected data block len: %d, got: %d", dataBlockLen, len(data))
	}

	for nonce < math.MaxUint64 {
		select {
		case <-exit:
			return hashInt, nonce, fmt.Errorf("calc interrupted")
		default:
		}
		hashInt, ok, err = CheckHash(data, nonce)
		if err != nil {
			return hashInt, nonce, err
		}
		if ok {
			break
		} else {
			nonce++
		}
	}
	if nonce == math.MaxUint64 {
		return hashInt, nonce, fmt.Errorf("cannot find proper hash")
	}

	return hashInt, nonce, nil
}

func CheckHash(data []byte, nonce uint64) (big.Int, bool, error) {
	var hashInt big.Int
	if len(data) != dataBlockLen {
		return hashInt, false, fmt.Errorf("expected data block len: %d, got: %d", dataBlockLen, len(data))
	}
	nonceByte := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceByte, nonce)
	store := append(data, nonceByte...)
	hash := sha256.Sum256(store)
	hashInt.SetBytes(hash[:])
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return hashInt, hashInt.Cmp(target) == -1, nil
}
