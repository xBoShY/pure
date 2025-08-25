package did

import (
	"fmt"
	"strings"
)

type DID struct {
	Network string
	Id      Id
}

func (did DID) String() string {
	return fmt.Sprintf("did:pure:%s:%s", did.Network, did.Id.String())
}

func DidDecode(didStr string) (DID, error) {
	var d DID
	var err error

	parts := strings.Split(didStr, ":")
	if len(parts) != 4 || parts[0] != "did" {
		return d, fmt.Errorf("invalid did %s", didStr)
	}

	if parts[1] != "pure" {
		return d, fmt.Errorf("unsupported did method %s", parts[1])
	}

	d.Id, err = DecodeId(parts[3])
	if err != nil {
		return d, err
	}

	d.Network = parts[2]

	return d, nil
}
