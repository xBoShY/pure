package did

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"fmt"
)

const (
	checksumLenBytes = 4
	hashLenBytes     = sha512.Size256
)

type Id [hashLenBytes]byte

func (id Id) String() string {
	// Compute the checksum
	checksumHash := sha512.Sum512_256(id[:])
	checksumLenBytes := checksumHash[hashLenBytes-checksumLenBytes:]

	// Append the checksum and encode as base32
	checksumAddress := append(id[:], checksumLenBytes...)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(checksumAddress)
}

func DecodeId(idStr string) (Id, error) {
	var id Id

	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(idStr)
	if err != nil {
		return id, err
	}

	if len(decoded) != len(id)+checksumLenBytes {
		return id, fmt.Errorf("decoded ID is the wrong length, should be %d bytes", hashLenBytes+checksumLenBytes)
	}

	// Split into address + checksum
	idBytes := decoded[:len(id)]
	checksumBytes := decoded[len(id):]

	// Compute the expected checksum
	checksumHash := sha512.Sum512_256(idBytes)
	expectedChecksumBytes := checksumHash[hashLenBytes-checksumLenBytes:]

	// Check the checksum
	if !bytes.Equal(expectedChecksumBytes, checksumBytes) {
		return id, fmt.Errorf("ID checksum is incorrect, did you copy the address correctly?")
	}

	// Checksum is good, copy address bytes into output
	copy(id[:], idBytes)

	// Check if address is canonical
	if id.String() != idStr {
		return id, fmt.Errorf("id %s is non-canonical", idStr)
	}

	return id, nil
}
