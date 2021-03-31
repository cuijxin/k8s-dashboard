package ecdsa

import (
	"crypto/elliptic"
	"testing"
)

func TestNewECDSACreator(t *testing.T) {
	keyFile := "cert.key"
	certFile := "cert.crt"
	creator := NewECDSACreator(keyFile, certFile, elliptic.P256())

	if creator == nil {
		t.Fatal("Expected creator not to be nil.")
	}
}

func TestEcdsaCreator_GetCertFileName(t *testing.T) {
	keyFile := "cert.key"
	certFile := "cert.crt"
	creator := NewECDSACreator(keyFile, certFile, elliptic.P256())

	if creator.GetCertFileName() != certFile {
		t.Fatalf("Expected cert file name to equal %s but go %s.", certFile, creator.GetCertFileName())
	}
}

func TestEcdsaCreator_GetKeyFileName(t *testing.T) {
	keyFile := "cert.key"
	certFile := "cert.crt"
	creator := NewECDSACreator(keyFile, certFile, elliptic.P256())

	if creator.GetKeyFileName() != keyFile {
		t.Fatalf("Expected cert key file name to equal %s but got %s.", keyFile, creator.GetKeyFileName())
	}
}
