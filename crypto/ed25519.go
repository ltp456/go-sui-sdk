package crypto

import (
	"crypto/ed25519"
	"fmt"
	"golang.org/x/crypto/sha3"
)

type KeyPair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	AuthKey    []byte
}

func (kp *KeyPair) Type() KeyType {
	return Ed25519Type
}

func (kp *KeyPair) Address() string {
	return fmt.Sprintf("0x%x", kp.AuthKey)
}

func (kp *KeyPair) Sign(data []byte) ([]byte, error) {
	var prefixBytes []byte
	signingMessage := append(prefixBytes, data...)
	return ed25519.Sign(kp.PrivateKey, signingMessage), nil
}

func NewKeyPairFromSeed(seed []byte) (*KeyPair, error) {
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	hash := sha3.New256()
	hash.Write([]byte{0x00})
	hash.Write(publicKey)
	hashSum := hash.Sum(nil)
	key := make([]byte, 20)
	copy(key, hashSum[:20])
	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		AuthKey:    key,
	}, nil
}

func (kp *KeyPair) Public() ed25519.PublicKey {
	return kp.PublicKey
}

func (kp *KeyPair) Verify(message, signature []byte) bool {
	return ed25519.Verify(kp.PublicKey, message, signature)
}
