package crypto

const SeedLength = 32
const SignatureLength = 64

type KeyType string

const (
	Ed25519Type   KeyType = "ED25519"
	Secp256k1Type KeyType = "SECP256K1TYPE"
)

func (k KeyType) String() string {
	return string(k)
}

type SigScheme byte

const (
	ED25519SigScheme   SigScheme = 0x00
	Secp256k1SigScheme SigScheme = 0x01
	BLS12381SigScheme  SigScheme = 0xff
)
