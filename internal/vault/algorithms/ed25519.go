package algorithms

type ed25519 struct {
}

var singleEd25519 = &ed25519{}

func Ed25519() *ed25519 {
	return singleEd25519
}

func (a *ed25519) Name() string {
	return "ed25519"
}

func (a *ed25519) Params() map[string]any {
	return map[string]any{
		"type": "ed25519",
	}
}

func (a *ed25519) Verify(message, signature []byte, key any) error {
	return nil
}

func (a *ed25519) KeyCheck(key any) error {
	return nil
}
