package did

import (
	"crypto/sha512"
	"encoding/base32"
)

type Id []byte

func (id Id) String() string {
	// Compute the checksum
	checksumHash := sha512.Sum512_256(id[:])
	checksumLenBytes := checksumHash[hashLenBytes-checksumLenBytes:]

	// Append the checksum and encode as base32
	checksumAddress := append(id[:], checksumLenBytes...)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(checksumAddress)
}
