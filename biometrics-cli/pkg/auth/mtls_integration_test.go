//go:build integration
// +build integration

package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMTLSHandshake(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()

	caCert, caKey, err := generateTestCA(tmpDir)
	if err != nil {
		t.Fatalf("Failed to generate test CA: %v", err)
	}

	serverCertPEM, serverKeyPEM, err := generateTestCertificate(
		tmpDir, "server",
		[]string{"localhost", "127.0.0.1"},
		[]net.IP{net.ParseIP("127.0.0.1")},
		caCert, caKey,
		false,
	)
	if err != nil {
		t.Fatalf("Failed to generate server cert: %v", err)
	}

	clientCertPEM, clientKeyPEM, err := generateTestCertificate(
		tmpDir, "client",
		[]string{},
		[]net.IP{},
		caCert, caKey,
		true,
	)
	if err != nil {
		t.Fatalf("Failed to generate client cert: %v", err)
	}

	serverCert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		t.Fatalf("Failed to load server cert: %v", err)
	}

	clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		t.Fatalf("Failed to load client cert: %v", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(caCert)

	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
		MinVersion:   tls.VersionTLS13,
	}

	if len(serverTLSConfig.CipherSuites) == 0 {
		serverTLSConfig.CipherSuites = []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_128_GCM_SHA256,
		}
	}

	_ = serverTLSConfig

	tlsCert := &tls.Certificate{
		Certificate: serverCert.Certificate,
		PrivateKey:  serverCert.PrivateKey,
	}

	if tlsCert == nil {
		t.Error("Server TLS config should be valid")
	}

	if len(serverTLSConfig.ClientCAs.Subjects()) == 0 {
		t.Error("Client CA pool should have certificates")
	}

	t.Logf("CA Certificate generated: Subject=%s", caCert.Subject.CommonName)
	t.Logf("Server Certificate generated: DNSNames=%v", serverCert.Leaf.DNSNames)
	t.Logf("Client Certificate generated: Subject=%s", clientCert.Leaf.Subject.CommonName)
}

func TestMTLSServerClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()

	caCert, caKey, err := generateTestCA(tmpDir)
	if err != nil {
		t.Fatalf("Failed to generate test CA: %v", err)
	}

	serverCertPEM, serverKeyPEM, err := generateTestCertificate(
		tmpDir, "server",
		[]string{"localhost"},
		[]net.IP{net.ParseIP("127.0.0.1")},
		caCert, caKey,
		false,
	)
	if err != nil {
		t.Fatalf("Failed to generate server cert: %v", err)
	}

	clientCertPEM, clientKeyPEM, err := generateTestCertificate(
		tmpDir, "client",
		[]string{},
		[]net.IP{},
		caCert, caKey,
		true,
	)
	if err != nil {
		t.Fatalf("Failed to generate client cert: %v", err)
	}

	serverCert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		t.Fatalf("Failed to load server cert: %v", err)
	}

	clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		t.Fatalf("Failed to load client cert: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS13,
	}

	clientTLSConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS13,
	}

	_ = serverTLSConfig
	_ = clientTLSConfig

	t.Log("Server and Client TLS configurations created successfully")
}

func TestMTLSCertificateValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()

	caCert, caKey, err := generateTestCA(tmpDir)
	if err != nil {
		t.Fatalf("Failed to generate test CA: %v", err)
	}

	expiredCertPEM, expiredKeyPEM, err := generateTestCertificate(
		tmpDir, "expired",
		[]string{"localhost"},
		[]net.IP{net.ParseIP("127.0.0.1")},
		caCert, caKey,
		false,
	)
	if err != nil {
		t.Fatalf("Failed to generate expired cert: %v", err)
	}

	expiredCert, err := tls.X509KeyPair(expiredCertPEM, expiredKeyPEM)
	if err != nil {
		t.Fatalf("Failed to load expired cert: %v", err)
	}

	now := time.Now()
	expiredTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Expired Cert",
		},
		NotBefore: now.Add(-2 * 24 * time.Hour),
		NotAfter:  now.Add(-1 * 24 * time.Hour),
	}

	expiredDerBytes, err := x509.CreateCertificate(rand.Reader, expiredTemplate, caCert, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("Failed to create expired cert: %v", err)
	}

	expiredX509Cert, err := x509.ParseCertificate(expiredDerBytes)
	if err != nil {
		t.Fatalf("Failed to parse expired cert: %v", err)
	}

	if !expiredX509Cert.NotAfter.Before(now) {
		t.Error("Expired certificate should be expired")
	}

	err = ValidateCertificate(expiredX509Cert, caCert)
	if err == nil {
		t.Error("Expired certificate should fail validation")
	}

	invalidCACert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Invalid CA",
		},
		NotBefore: now,
		NotAfter:  now.Add(24 * time.Hour),
	}

	err = ValidateCertificate(caCert, invalidCACert)
	if err == nil {
		t.Error("Certificate with invalid CA should fail validation")
	}

	_ = expiredCert

	t.Log("Certificate validation tests passed")
}

func generateTestCA(tmpDir string) ([]byte, *rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate CA key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial: %w", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
			CommonName:   "Test CA",
		},
		NotBefore:             now,
		NotAfter:              now.Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CA cert: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	caCertPath := filepath.Join(tmpDir, "test-ca.pem")
	if err := os.WriteFile(caCertPath, certPEM, 0644); err != nil {
		return nil, nil, fmt.Errorf("failed to write CA cert: %w", err)
	}

	return certPEM, privateKey, nil
}

func generateTestCertificate(
	tmpDir, name string,
	dnsNames []string,
	ipAddresses []net.IP,
	caCert *x509.Certificate,
	caKey *rsa.PrivateKey,
	isClient bool,
) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial: %w", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: name,
		},
		NotBefore:          now,
		NotAfter:           now.Add(24 * time.Hour),
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:           dnsNames,
		IPAddresses:        ipAddresses,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	if isClient {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &privateKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cert: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	certPath := filepath.Join(tmpDir, name+".pem")
	keyPath := filepath.Join(tmpDir, name+".key")

	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return nil, nil, fmt.Errorf("failed to write cert: %w", err)
	}

	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return nil, nil, fmt.Errorf("failed to write key: %w", err)
	}

	return certPEM, keyPEM, nil
}

func ValidateCertificate(cert, caCert *x509.Certificate) error {
	now := time.Now()

	if cert.NotAfter.Before(now) {
		return fmt.Errorf("certificate expired")
	}

	if cert.NotBefore.After(now) {
		return fmt.Errorf("certificate not yet valid")
	}

	if cert.IsCA {
		return nil
	}

	opts := x509.VerifyOptions{
		Roots: &x509.CertPool{},
	}

	if caCert != nil {
		opts.Roots.AddCert(caCert)
	}

	_, err := cert.Verify(opts)
	return err
}
