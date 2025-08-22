package did

import (
	"crypto/ed25519"
	"crypto/sha512"
	"errors"
)

const (
	checksumLenBytes = 4
	hashLenBytes     = sha512.Size256
)

type Object [hashLenBytes]byte

func (o Object) Id() Id {
	return []byte(o[:])
}

func NewObject(pub ed25519.PublicKey) (Object, error) {
	var o Object
	n := copy(o[:], pub)
	if n != ed25519.PublicKeySize {
		return o, errors.New("generated public key is the wrong size")
	}

	return o, nil
}

func (o Object) DID(network string) (DID, error) {
	return DID{
		Network: network,
		Id:      o.Id(),
	}, nil
}
