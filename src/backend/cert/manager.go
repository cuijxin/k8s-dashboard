package cert

import (
	"crypto/tls"
	certapi "github.com/cuijxin/k8s-dashboard/src/backend/cert/api"
	"log"
	"os"
)

// Manager is used to implement cert/api/types.Manager interface. See Manager for more information.
type Manager struct {
	creator certapi.Creator
	certDir string
}

// GetCertificates implements Manager interface. See Manager for more information.
func (m *Manager) GetCertificates() (tls.Certificate, error) {
	if m.keyFileExists() && m.certFileExists() {
		log.Println("Certificates already exists. Returning.")
		return tls.LoadX509KeyPair(
			m.path(m.creator.GetCertFileName()),
			m.path(m.creator.GetKeyFileName()),
		)
	}

	key := m.creator.GenerateKey()
	cert := m.creator.GenerateCertificate(key)
	log.Println("Successfully created certificates")
	keyPEM, certPEM, err := m.creator.KeyCertPEMBytes(key, cert)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certPEM, keyPEM)
}

func (m *Manager) keyFileExists() bool {
	return m.exists(m.path(m.creator.GetKeyFileName()))
}

func (m *Manager) certFileExists() bool {
	return m.exists(m.path(m.creator.GetCertFileName()))
}

func (m *Manager) path(certFile string) string {
	return m.certDir + string(os.PathSeparator) + certFile
}

func (m *Manager) exists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

// NewCertManager creates Manager object.
func NewCertManager(creator certapi.Creator, certDir string) certapi.Manager {
	return &Manager{creator: creator, certDir: certDir}
}
