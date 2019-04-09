package model

import "crypto/rsa"

type Keys struct {
	PublicKeyFile  string
	PrivateKeyFile string
	PublicKey      *rsa.PublicKey
	PrivateKey     *rsa.PrivateKey
}
