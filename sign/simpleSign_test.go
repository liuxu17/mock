package sign

import (
	"fmt"
	"github.com/irisnet/irishub/crypto/keys/hd"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tyler-smith/go-bip39"
	"log"
	"testing"
)

func TestSign(t *testing.T){
	mnemonic := "tube lonely pause spring gym veteran know want grid tired taxi such same mesh charge orient bracket ozone concert once good quick dry boss"
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		log.Fatal(err)
	}

	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hd.FullFundraiserPath)
	fmt.Println(derivedPriv)
	pubk := secp256k1.PrivKeySecp256k1(derivedPriv).PubKey()

	pubk.Bytes()
	fmt.Println(pubk.Bytes())

}
