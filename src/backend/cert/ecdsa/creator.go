package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	certapi "github.com/cuijxin/k8s-dashboard/src/backend/cert/api"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Implements certificate Creator interface. See Creator for more information.
type ecdsaCreator struct {
	keyFile  string
	certFile string
	curve    elliptic.Curve
}

// GenerateKey implements certificate Creator interface. See Creator for more
// information.
func (c *ecdsaCreator) GenerateKey() interface{} {
	key, err := ecdsa.GenerateKey(c.curve, rand.Reader)
	if err != nil {
		log.Fatalf("[ECDSAManager] Failed to generate certificate key: %s", err)
	}

	return key
}

// GenerateCertificate implements certificate Creator interface. See Creator for
// more information.
func (c *ecdsaCreator) GenerateCertificate(key interface{}) []byte {
	ecdsaKey := c.getKey(key)
	pod := c.getDashboardPod()

	notBefore := time.Now()
	validFor, _ := time.ParseDuration("8760h")
	notAfter := notBefore.Add(validFor)

	template := x509.Certificate{
		SerialNumber: c.generateSerialNumber(),
		NotAfter:     notAfter,
		NotBefore:    notBefore,
	}

	if len(pod.Name) > 0 && len(pod.Namespace) > 0 {
		podDomainName := pod.Name + "." + pod.Namespace
		template.Subject = pkix.Name{CommonName: podDomainName}
		template.Issuer = pkix.Name{CommonName: podDomainName}
		template.DNSNames = []string{podDomainName}
	}

	if len(pod.Status.PodIP) > 0 {
		template.IPAddresses = []net.IP{net.ParseIP(pod.Status.PodIP)}
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &ecdsaKey.PublicKey, ecdsaKey)
	if err != nil {
		log.Fatalf("[ECDSAManager] Failed to create certificate: %s", err)
	}

	return certBytes
}

// StoreCertificates implements certificate Creator interface. See Creator for more information.
func (c *ecdsaCreator) StoreCertificates(path string, key interface{}, certBytes []byte) {
	keyPEM, certPEM, err := c.KeyCertPEMBytes(key, certBytes)
	if err != nil {
		log.Fatalf("[ECDSAManager] Failed to marshal cert/key pair: %v", err)
	}
	if err := ioutil.WriteFile(path+string(os.PathSeparator)+c.GetCertFileName(), certPEM, os.FileMode(0644)); err != nil {
		log.Fatalf("[ECDSAManager] Failed to open %s for writing: %s", c.GetCertFileName(), err)
	}
	if err := ioutil.WriteFile(path+string(os.PathSeparator)+c.GetKeyFileName(), keyPEM, os.FileMode(0600)); err != nil {
		log.Fatalf("[ECDSAManager] Failed to open %s for writing: %s", c.GetKeyFileName(), err)
	}
}

func (c *ecdsaCreator) KeyCertPEMBytes(key interface{}, certBytes []byte) ([]byte, []byte, error) {
	marshaledKey, err := x509.MarshalECPrivateKey(c.getKey(key))
	if err != nil {
		return nil, nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: marshaledKey})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	return keyPEM, certPEM, nil
}

// GetKeyFileName implements certificate Creator interface. See Creator for more information.
func (c *ecdsaCreator) GetKeyFileName() string {
	return c.keyFile
}

// GetCertFileName implements certificate Creator interface. See Creator for more information.
func (c *ecdsaCreator) GetCertFileName() string {
	return c.certFile
}

func (c *ecdsaCreator) getKey(key interface{}) *ecdsa.PrivateKey {
	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		log.Fatal("[ECDSAManager] Key should be an instance of *ecdsa.PrivateKey")
	}

	return ecdsaKey
}

func (c *ecdsaCreator) generateSerialNumber() *big.Int {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("[ECDSAManager] Failed to generate serial number: %s", err)
	}

	return serialNumber
}

func (c *ecdsaCreator) getDashboardPod() *corev1.Pod {
	// These variables are populated by kubernetes downward API when using in-cluster
	// config
	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	podIP := os.Getenv("POD_IP")

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: podNamespace,
		},
		Status: corev1.PodStatus{
			PodIP: podIP,
		},
	}
}

func (c *ecdsaCreator) init() {
	if len(c.certFile) == 0 {
		c.certFile = certapi.DashboardCertName
	}

	if len(c.keyFile) == 0 {
		c.keyFile = certapi.DashboardKeyName
	}
}

// NewECDSACreator creates ECDSACreator instance
func NewECDSACreator(keyFile, certFile string, curve elliptic.Curve) certapi.Creator {
	creator := &ecdsaCreator{
		curve:    curve,
		keyFile:  keyFile,
		certFile: certFile,
	}

	creator.init()
	return creator
}
