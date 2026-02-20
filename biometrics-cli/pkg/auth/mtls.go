package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	ErrInvalidCertificate     = errors.New("invalid certificate")
	ErrCertificateExpired     = errors.New("certificate expired")
	ErrCertificateNotYetValid = errors.New("certificate not yet valid")
	ErrInvalidCA              = errors.New("invalid certificate authority")
	ErrMTLSHandshakeFailed    = errors.New("mTLS handshake failed")
)

type MTLSConfig struct {
	CAPath           string
	ServerCertPath   string
	ServerKeyPath    string
	ClientCAPath     string
	MinVersion       uint16
	MaxVersion       uint16
	CipherSuites     []uint16
	ClientAuth       tls.ClientAuthType
	RotationInterval time.Duration
	AutoRotate       bool
}

type MTLSManager struct {
	config       *MTLSConfig
	caCert       *x509.Certificate
	caKey        *rsa.PrivateKey
	serverCert   *tls.Certificate
	clientCAPool *x509.CertPool
	mu           sync.RWMutex
	lastRotation time.Time
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

type CertificateInfo struct {
	SerialNumber *big.Int
	Subject      pkix.Name
	Issuer       pkix.Name
	NotBefore    time.Time
	NotAfter     time.Time
	DNSNames     []string
	IPAddresses  []net.IP
	KeyUsage     x509.KeyUsage
	ExtKeyUsage  []x509.ExtKeyUsage
	IsCA         bool
	SignatureAlg x509.SignatureAlgorithm
	PublicKeyAlg string
}

func DefaultMTLSConfig() *MTLSConfig {
	return &MTLSConfig{
		CAPath:           "/tmp/biometrics/mtls/ca.pem",
		ServerCertPath:   "/tmp/biometrics/mtls/server.pem",
		ServerKeyPath:    "/tmp/biometrics/mtls/server.key",
		ClientCAPath:     "/tmp/biometrics/mtls/client-ca.pem",
		MinVersion:       tls.VersionTLS13,
		MaxVersion:       tls.VersionTLS13,
		CipherSuites:     nil,
		ClientAuth:       tls.RequireAndVerifyClientCert,
		RotationInterval: 24 * time.Hour,
		AutoRotate:       false,
	}
}

func NewMTLSManager(config *MTLSConfig) (*MTLSManager, error) {
	if config == nil {
		config = DefaultMTLSConfig()
	}

	manager := &MTLSManager{
		config:   config,
		stopChan: make(chan struct{}),
	}

	if err := manager.loadOrGenerateCA(); err != nil {
		return nil, fmt.Errorf("failed to load CA: %w", err)
	}

	if err := manager.loadServerCertificate(); err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	if err := manager.loadClientCA(); err != nil {
		return nil, fmt.Errorf("failed to load client CA: %w", err)
	}

	if config.AutoRotate {
		manager.startAutoRotation()
	}

	return manager, nil
}

func (m *MTLSManager) loadOrGenerateCA() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if data, err := os.ReadFile(m.config.CAPath); err == nil {
		block, _ := pem.Decode(data)
		if block == nil {
			return ErrInvalidCA
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse CA certificate: %w", err)
		}

		m.caCert = cert

		keyPath := filepath.Join(filepath.Dir(m.config.CAPath), "ca.key")
		if keyData, err := os.ReadFile(keyPath); err == nil {
			block, _ := pem.Decode(keyData)
			if block == nil {
				return ErrInvalidCA
			}

			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse CA key: %w", err)
			}

			m.caKey = key
			return nil
		}

		return ErrInvalidCA
	}

	return m.generateCA()
}

func (m *MTLSManager) generateCA() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate CA key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"BIOMETRICS"},
			Country:      []string{"DE"},
			Province:     []string{"Berlin"},
			Locality:     []string{"Berlin"},
			CommonName:   "BIOMETRICS Root CA",
		},
		NotBefore:             now,
		NotAfter:              now.Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		SignatureAlgorithm:    x509.SHA256WithRSA,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	m.caCert = cert
	m.caKey = privateKey

	if err := m.saveCertificate(cert, m.config.CAPath); err != nil {
		return fmt.Errorf("failed to save CA certificate: %w", err)
	}

	keyPath := filepath.Join(filepath.Dir(m.config.CAPath), "ca.key")
	if err := m.savePrivateKey(privateKey, keyPath); err != nil {
		return fmt.Errorf("failed to save CA key: %w", err)
	}

	return nil
}

func (m *MTLSManager) loadServerCertificate() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := os.Stat(m.config.ServerCertPath); os.IsNotExist(err) {
		return m.generateServerCertificate()
	}

	cert, err := tls.LoadX509KeyPair(m.config.ServerCertPath, m.config.ServerKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load server certificate: %w", err)
	}

	m.serverCert = &cert
	return nil
}

func (m *MTLSManager) generateServerCertificate() error {
	if m.caCert == nil || m.caKey == nil {
		return ErrInvalidCA
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate server key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"BIOMETRICS"},
			CommonName:   "BIOMETRICS Server",
		},
		NotBefore:          now,
		NotAfter:           now.Add(90 * 24 * time.Hour),
		KeyUsage:           x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:           []string{"localhost", "biometrics.local"},
		IPAddresses:        []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, m.caCert, &privateKey.PublicKey, m.caKey)
	if err != nil {
		return fmt.Errorf("failed to create server certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return fmt.Errorf("failed to parse server certificate: %w", err)
	}

	tlsCert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privateKey,
		Leaf:        cert,
	}

	m.serverCert = &tlsCert

	if err := m.saveCertificate(cert, m.config.ServerCertPath); err != nil {
		return fmt.Errorf("failed to save server certificate: %w", err)
	}

	if err := m.savePrivateKey(privateKey, m.config.ServerKeyPath); err != nil {
		return fmt.Errorf("failed to save server key: %w", err)
	}

	return nil
}

func (m *MTLSManager) loadClientCA() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clientCAPool = x509.NewCertPool()
	m.clientCAPool.AddCert(m.caCert)

	if m.config.ClientCAPath != "" {
		if data, err := os.ReadFile(m.config.ClientCAPath); err == nil {
			if !m.clientCAPool.AppendCertsFromPEM(data) {
				return ErrInvalidCA
			}
		}
	}

	return nil
}

func (m *MTLSManager) GetTLSConfig() (*tls.Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.serverCert == nil {
		return nil, ErrInvalidCertificate
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{*m.serverCert},
		ClientAuth:   m.config.ClientAuth,
		ClientCAs:    m.clientCAPool,
		MinVersion:   m.config.MinVersion,
		MaxVersion:   m.config.MaxVersion,
		CipherSuites: m.config.CipherSuites,
	}

	if len(config.CipherSuites) == 0 {
		config.CipherSuites = []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_128_GCM_SHA256,
		}
	}

	return config, nil
}

func (m *MTLSManager) ValidateClientCertificate(cert *x509.Certificate) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.clientCAPool == nil {
		return ErrInvalidCA
	}

	now := time.Now()
	if cert.NotAfter.Before(now) {
		return ErrCertificateExpired
	}
	if cert.NotBefore.After(now) {
		return ErrCertificateNotYetValid
	}

	opts := x509.VerifyOptions{
		Roots:     m.clientCAPool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("certificate verification failed: %w", err)
	}

	return nil
}

func (m *MTLSManager) GenerateClientCertificate(clientID string, validity time.Duration) (*tls.Certificate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.caCert == nil || m.caKey == nil {
		return nil, ErrInvalidCA
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"BIOMETRICS"},
			CommonName:   fmt.Sprintf("Client: %s", clientID),
		},
		NotBefore:          now,
		NotAfter:           now.Add(validity),
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, m.caCert, &privateKey.PublicKey, m.caKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create client certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse client certificate: %w", err)
	}

	return &tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privateKey,
		Leaf:        cert,
	}, nil
}

func (m *MTLSManager) GetCertificateInfo(certPath string) (*CertificateInfo, error) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidCertificate
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return &CertificateInfo{
		SerialNumber: cert.SerialNumber,
		Subject:      cert.Subject,
		Issuer:       cert.Issuer,
		NotBefore:    cert.NotBefore,
		NotAfter:     cert.NotAfter,
		DNSNames:     cert.DNSNames,
		IPAddresses:  cert.IPAddresses,
		KeyUsage:     cert.KeyUsage,
		ExtKeyUsage:  cert.ExtKeyUsage,
		IsCA:         cert.IsCA,
		SignatureAlg: cert.SignatureAlgorithm,
		PublicKeyAlg: cert.PublicKeyAlgorithm.String(),
	}, nil
}

func (m *MTLSManager) RotateCertificates() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.generateServerCertificate(); err != nil {
		return fmt.Errorf("failed to generate new server certificate: %w", err)
	}

	m.lastRotation = time.Now()
	return nil
}

func (m *MTLSManager) startAutoRotation() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		ticker := time.NewTicker(m.config.RotationInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := m.RotateCertificates(); err != nil {
					fmt.Printf("Failed to rotate certificates: %v\n", err)
				}
			case <-m.stopChan:
				return
			}
		}
	}()
}

func (m *MTLSManager) Stop() {
	close(m.stopChan)
	m.wg.Wait()
}

func (m *MTLSManager) saveCertificate(cert *x509.Certificate, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})

	return os.WriteFile(path, certPEM, 0644)
}

func (m *MTLSManager) savePrivateKey(key *rsa.PrivateKey, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return os.WriteFile(path, keyPEM, 0600)
}

func (m *MTLSManager) GetServerCertificate() *tls.Certificate {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.serverCert
}

func (m *MTLSManager) GetCACertificate() *x509.Certificate {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.caCert
}

func IsCertificateValid(cert *x509.Certificate) bool {
	now := time.Now()
	return !now.Before(cert.NotBefore) && !now.After(cert.NotAfter)
}

func DaysUntilExpiry(cert *x509.Certificate) int {
	return int(time.Until(cert.NotAfter).Hours() / 24)
}

func (m *MTLSManager) NeedsRotation() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.serverCert == nil || len(m.serverCert.Certificate) == 0 {
		return true
	}

	cert, err := x509.ParseCertificate(m.serverCert.Certificate[0])
	if err != nil {
		return true
	}

	return DaysUntilExpiry(cert) < 7
}
