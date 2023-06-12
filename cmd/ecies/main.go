package main

import (
	"fmt"
	"os"

	ecies "github.com/ecies/go/v2"
)

func main() {
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
	if err := os.WriteFile("pub-key.pem", []byte(pub.Hex(false)), 0644); err != nil {
		panic(err)
	}
}
