package util

import (
    "crypto/sha256"
)

func GetHash(data []byte) [32]byte {
    return sha256.Sum256(data)
}