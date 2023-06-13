package main

import (
	"fmt"
	"os"

	hpke "github.com/jedisct1/go-hpke-compact"
	ecies "github.com/yyewolf/go-ecies/v2"
)

func eciesGenKey() {
	fmt.Println("Generating ECIES keypair...")
	priv, err := ecies.GenerateKey()
	if err != nil {
		panic(err)
	}
	pub := priv.PublicKey

	fmt.Println("Storing ECIES keypair into priv-key.pem and pub-key.pem...")
	if err := os.WriteFile("priv-key.pem", []byte(priv.Hex()), 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile("pub-key.pem", []byte(pub.Hex()), 0644); err != nil {
		panic(err)
	}
}

func hpkeGenKey() {
	fmt.Println("Generating ECIES keypair...")
	suite, err := hpke.NewSuite(hpke.KemX25519HkdfSha256, hpke.KdfHkdfSha256, hpke.AeadChaCha20Poly1305)
	if err != nil {
		panic(err)
	}
	clientKp, err := suite.GenerateKeyPair()
	if err != nil {
		panic(err)
	}
	serverKp, err := suite.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	fmt.Println("Storing HPKE keypair into client-priv-key.pem, client-pub-key.pem, server-priv-key.pem, server-pub-key.pem...")
	if err := os.WriteFile("client-priv-key.pem", clientKp.SecretKey, 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile("client-pub-key.pem", clientKp.PublicKey, 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile("server-priv-key.pem", serverKp.SecretKey, 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile("server-pub-key.pem", serverKp.PublicKey, 0644); err != nil {
		panic(err)
	}
}
