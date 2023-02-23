package crypto

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestNewKeyPairFromSeed(t *testing.T) {
	seed := "0efd0145c9854b3189b20201e93b0fa91bd68b95936363f172846150fca902d7"
	hexSeed, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	keyPair, err := NewKeyPairFromSeed(hexSeed)
	if err != nil {
		panic(err)
	}
	fmt.Printf("address: %v\n", keyPair.Address())

	signature, err := keyPair.Sign([]byte("wo he ni"))
	if err != nil {
		panic(err)
	}

	verify := keyPair.Verify([]byte("wo he ni"), signature)
	fmt.Printf("verify result: %v\n", verify)
	verify1 := keyPair.Verify([]byte("wo he nidd"), signature)
	fmt.Printf("verify result01: %v\n", verify1)

}
