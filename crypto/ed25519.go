package crypto

import (
	"crypto/ed25519"
	"fmt"
	"golang.org/x/crypto/blake2b"
)

type KeyPair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	AuthKey    [32]byte
}

func (kp *KeyPair) Type() KeyType {
	return Ed25519Type
}

func (kp *KeyPair) Address() string {
	return fmt.Sprintf("0x%x", kp.AuthKey)
}

func (kp *KeyPair) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(kp.PrivateKey, data), nil
}

func NewKeyPairFromSeed(seed []byte) (*KeyPair, error) {
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)
	data := make([]byte, 0)
	data = append(data, []byte{0x00}...)
	data = append(data, publicKey...)
	hash := blake2b.Sum256(data)
	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		AuthKey:    hash,
	}, nil
}

func (kp *KeyPair) Public() ed25519.PublicKey {
	return kp.PublicKey
}

func (kp *KeyPair) Verify(message, signature []byte) bool {
	return ed25519.Verify(kp.PublicKey, message, signature)
}
