package did

import (
	"crypto/ed25519"
	"errors"
)

type Object = Id

func (o Object) Id() Id {
	return o
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

func (o Object) CanonicalPublicKey() (ed25519.PublicKey, error) {
	var p ed25519.PublicKey = make(ed25519.PublicKey, ed25519.PublicKeySize)
	copy(p, o[:])

	return p, nil
}
