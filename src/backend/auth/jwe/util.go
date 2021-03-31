package jwe

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
)

// Credits to David W. https://stackoverflow.com/a/44688503

// ExportRSAKeyOrDie exports rsa key object to a private/public strings. In case
// of fail panic is called.
func ExportRSAKeyOrDie(privKey *rsa.PrivateKey) (priv, pub string) {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privkeyPems := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)

	priv = string(privkeyPems)

	pubkeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		panic(err)
	}

	pubkeyPems := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	pub = string(pubkeyPems)
	return
}

// ParseRSAKey parses private/public key strings and returns rsa key object or error.
func ParseRSAKey(privStr, pubStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privStr))
	if block == nil {
		return nil, errors.NewInvalid("Failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	block, _ = pem.Decode([]byte(pubStr))
	if block == nil {
		return nil, errors.NewInvalid("Failed to parse PEM block containing the key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.NewInvalid("Failed to parse public key")
	}

	priv.PublicKey = *pub
	return priv, nil
}
