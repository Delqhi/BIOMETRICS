package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

type CertificateBundle struct {
	CA         *x509.Certificate
	CAKey      *rsa.PrivateKey
	ServerCert *x509.Certificate
	ServerKey  *rsa.PrivateKey
	ClientCert *x509.Certificate
	ClientKey  *rsa.PrivateKey
}

type CertificateGenerator struct {
	organization string
	country      string
	validity     time.Duration
	keySize      int
}

func NewCertificateGenerator(org, country string, validity time.Duration) *CertificateGenerator {
	return &CertificateGenerator{
		organization: org,
		country:      country,
		validity:     validity,
		keySize:      4096,
	}
}

func (g *CertificateGenerator) GenerateCA() (*CertificateBundle, error) {
	caKey, err := rsa.GenerateKey(rand.Reader, g.keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate CA key: %w", err)
	}

	caCert, err := g.createCACertificate(caKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create CA certificate: %w", err)
	}

	return &CertificateBundle{
		CA:    caCert,
		CAKey: caKey,
	}, nil
}

func (g *CertificateGenerator) createCACertificate(caKey *rsa.PrivateKey) (*x509.Certificate, error) {
	serialNumber, err := generateSerialNumber()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{g.organization},
			Country:      []string{g.country},
			CommonName:   fmt.Sprintf("%s Root CA", g.organization),
		},
		NotBefore:             now,
		NotAfter:              now.Add(g.validity),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create CA certificate: %w", err)
	}

	return x509.ParseCertificate(derBytes)
}

func (g *CertificateGenerator) GenerateServerCertificate(caBundle *CertificateBundle, dnsNames []string, ips []net.IP) (*CertificateBundle, error) {
	if caBundle == nil || caBundle.CA == nil || caBundle.CAKey == nil {
		return nil, fmt.Errorf("invalid CA bundle")
	}

	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate server key: %w", err)
	}

	serverCert, err := g.createServerCertificate(serverKey, caBundle.CA, caBundle.CAKey, dnsNames, ips)
	if err != nil {
		return nil, fmt.Errorf("failed to create server certificate: %w", err)
	}

	return &CertificateBundle{
		CA:         caBundle.CA,
		CAKey:      caBundle.CAKey,
		ServerCert: serverCert,
		ServerKey:  serverKey,
	}, nil
}

func (g *CertificateGenerator) createServerCertificate(serverKey *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey, dnsNames []string, ips []net.IP) (*x509.Certificate, error) {
	serialNumber, err := generateSerialNumber()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{g.organization},
			CommonName:   fmt.Sprintf("%s Server", g.organization),
		},
		NotBefore:          now,
		NotAfter:           now.Add(90 * 24 * time.Hour),
		KeyUsage:           x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:           dnsNames,
		IPAddresses:        ips,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &serverKey.PublicKey, caKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create server certificate: %w", err)
	}

	return x509.ParseCertificate(derBytes)
}

func (g *CertificateGenerator) GenerateClientCertificate(caBundle *CertificateBundle, clientID string) (*CertificateBundle, error) {
	if caBundle == nil || caBundle.CA == nil || caBundle.CAKey == nil {
		return nil, fmt.Errorf("invalid CA bundle")
	}

	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client key: %w", err)
	}

	clientCert, err := g.createClientCertificate(clientKey, caBundle.CA, caBundle.CAKey, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to create client certificate: %w", err)
	}

	return &CertificateBundle{
		CA:         caBundle.CA,
		CAKey:      caBundle.CAKey,
		ClientCert: clientCert,
		ClientKey:  clientKey,
	}, nil
}

func (g *CertificateGenerator) createClientCertificate(clientKey *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey, clientID string) (*x509.Certificate, error) {
	serialNumber, err := generateSerialNumber()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{g.organization},
			CommonName:   fmt.Sprintf("Client: %s", clientID),
		},
		NotBefore:          now,
		NotAfter:           now.Add(365 * 24 * time.Hour),
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &clientKey.PublicKey, caKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create client certificate: %w", err)
	}

	return x509.ParseCertificate(derBytes)
}

func SaveCertificateBundle(bundle *CertificateBundle, dir string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if bundle.CA != nil {
		if err := saveCertificate(bundle.CA, filepath.Join(dir, "ca.pem")); err != nil {
			return err
		}
	}

	if bundle.CAKey != nil {
		if err := savePrivateKey(bundle.CAKey, filepath.Join(dir, "ca.key")); err != nil {
			return err
		}
	}

	if bundle.ServerCert != nil {
		if err := saveCertificate(bundle.ServerCert, filepath.Join(dir, "server.pem")); err != nil {
			return err
		}
	}

	if bundle.ServerKey != nil {
		if err := savePrivateKey(bundle.ServerKey, filepath.Join(dir, "server.key")); err != nil {
			return err
		}
	}

	if bundle.ClientCert != nil {
		if err := saveCertificate(bundle.ClientCert, filepath.Join(dir, "client.pem")); err != nil {
			return err
		}
	}

	if bundle.ClientKey != nil {
		if err := savePrivateKey(bundle.ClientKey, filepath.Join(dir, "client.key")); err != nil {
			return err
		}
	}

	return nil
}

func LoadCertificateBundle(dir string) (*CertificateBundle, error) {
	bundle := &CertificateBundle{}

	if caCert, err := loadCertificate(filepath.Join(dir, "ca.pem")); err == nil {
		bundle.CA = caCert
	}

	if caKey, err := loadPrivateKey(filepath.Join(dir, "ca.key")); err == nil {
		bundle.CAKey = caKey
	}

	if serverCert, err := loadCertificate(filepath.Join(dir, "server.pem")); err == nil {
		bundle.ServerCert = serverCert
	}

	if serverKey, err := loadPrivateKey(filepath.Join(dir, "server.key")); err == nil {
		bundle.ServerKey = serverKey
	}

	if clientCert, err := loadCertificate(filepath.Join(dir, "client.pem")); err == nil {
		bundle.ClientCert = clientCert
	}

	if clientKey, err := loadPrivateKey(filepath.Join(dir, "client.key")); err == nil {
		bundle.ClientKey = clientKey
	}

	return bundle, nil
}

func saveCertificate(cert *x509.Certificate, path string) error {
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})

	return os.WriteFile(path, certPEM, 0644)
}

func savePrivateKey(key *rsa.PrivateKey, path string) error {
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return os.WriteFile(path, keyPEM, 0600)
}

func loadCertificate(path string) (*x509.Certificate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode certificate")
	}

	return x509.ParseCertificate(block.Bytes)
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func ExportCertificateToPEM(cert *x509.Certificate) ([]byte, error) {
	var buf bytes.Buffer
	err := pem.Encode(&buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ExportPrivateKeyToPEM(key *rsa.PrivateKey) ([]byte, error) {
	var buf bytes.Buffer
	err := pem.Encode(&buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ImportCertificateFromPEM(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}

func ImportPrivateKeyFromPEM(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func VerifyCertificateChain(cert *x509.Certificate, roots *x509.CertPool, intermediates []*x509.Certificate) error {
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
	}

	for _, intermediate := range intermediates {
		opts.Intermediates.AddCert(intermediate)
	}

	_, err := cert.Verify(opts)
	return err
}

func CheckCertificateExpiry(cert *x509.Certificate, warningDays int) (bool, error) {
	if cert == nil {
		return false, fmt.Errorf("certificate is nil")
	}

	now := time.Now()
	if now.After(cert.NotAfter) {
		return true, fmt.Errorf("certificate expired on %s", cert.NotAfter.Format(time.RFC3339))
	}

	if now.Before(cert.NotBefore) {
		return false, fmt.Errorf("certificate not yet valid, valid from %s", cert.NotBefore.Format(time.RFC3339))
	}

	daysUntilExpiry := int(time.Until(cert.NotAfter).Hours() / 24)
	if daysUntilExpiry < warningDays {
		return true, fmt.Errorf("certificate expires in %d days", daysUntilExpiry)
	}

	return false, nil
}

func GetCertificateFingerprint(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}

	hash := sha256.Sum256(cert.Raw)
	return fmt.Sprintf("%X", hash)
}

func CertificatesEqual(cert1, cert2 *x509.Certificate) bool {
	if cert1 == nil || cert2 == nil {
		return false
	}
	return bytes.Equal(cert1.Raw, cert2.Raw)
}

func generateSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}
	return serialNumber, nil
}
