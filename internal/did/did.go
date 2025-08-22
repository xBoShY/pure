package did

import "fmt"

type DID struct {
	Network string
	Id      Id
}

func (did DID) String() string {
	return fmt.Sprintf("did:pure:%s:%s", did.Network, did.Id.String())
}
