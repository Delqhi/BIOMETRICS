package cert

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrInvalidKeyPair      = errors.New("invalid key pair")
	ErrCertificateExpired  = errors.New("certificate has expired")
	ErrCertificateNotValid = errors.New("certificate is not yet valid")
	ErrInvalidCertificate  = errors.New("invalid certificate")
	ErrKeyMismatch         = errors.New("private key does not match certificate public key")
)

type KeyType int

const (
	KeyTypeRSA KeyType = iota
	KeyTypeECDSA
	KeyTypeEd25519
)

type CertificateConfig struct {
	Subject      pkix.Name
	DNSNames     []string
	IPAddresses  []net.IP
	KeyUsage     x509.KeyUsage
	ExtKeyUsage  []x509.ExtKeyUsage
	NotBefore    time.Time
	NotAfter     time.Time
	KeyType      KeyType
	RSAKeySize   int
	ECDSACurve   elliptic.Curve
	SerialNumber *big.Int
}

type Certificate struct {
	Certificate *x509.Certificate
	PrivateKey  interface{}
	PublicKey   interface{}
	CertPEM     []byte
	KeyPEM      []byte
}

type CAConfig struct {
	Subject      pkix.Name
	KeyType      KeyType
	RSAKeySize   int
	ECDSACurve   elliptic.Curve
	NotBefore    time.Time
	NotAfter     time.Time
	SerialNumber *big.Int
}

type ServerCertConfig struct {
	DNSNames    []string
	IPAddresses []net.IP
	NotBefore   time.Time
	NotAfter    time.Time
	KeyType     KeyType
	RSAKeySize  int
	ECDSACurve  elliptic.Curve
}

type ClientCertConfig struct {
	Subject    pkix.Name
	NotBefore  time.Time
	NotAfter   time.Time
	KeyType    KeyType
	RSAKeySize int
	ECDSACurve elliptic.Curve
}

func DefaultCAConfig() *CAConfig {
	return &CAConfig{
		Subject: pkix.Name{
			Organization: []string{"BIOMETRICS CLI"},
			Country:      []string{"US"},
			Province:     []string{"California"},
			Locality:     []string{"San Francisco"},
			CommonName:   "BIOMETRICS CLI Root CA",
		},
		KeyType:    KeyTypeECDSA,
		ECDSACurve: elliptic.P256(),
		NotBefore:  time.Now(),
		NotAfter:   time.Now().AddDate(10, 0, 0),
	}
}

func DefaultServerCertConfig() *ServerCertConfig {
	return &ServerCertConfig{
		DNSNames:    []string{"localhost"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0),
		KeyType:     KeyTypeECDSA,
		ECDSACurve:  elliptic.P256(),
	}
}

func DefaultClientCertConfig() *ClientCertConfig {
	return &ClientCertConfig{
		Subject: pkix.Name{
			Organization: []string{"BIOMETRICS CLI"},
			CommonName:   "BIOMETRICS CLI Client",
		},
		NotBefore:  time.Now(),
		NotAfter:   time.Now().AddDate(1, 0, 0),
		KeyType:    KeyTypeECDSA,
		ECDSACurve: elliptic.P256(),
	}
}

func generatePrivateKey(keyType KeyType, rsaKeySize int, ecdsaCurve elliptic.Curve) (interface{}, error) {
	switch keyType {
	case KeyTypeRSA:
		if rsaKeySize == 0 {
			rsaKeySize = 2048
		}
		return rsa.GenerateKey(rand.Reader, rsaKeySize)
	case KeyTypeECDSA:
		if ecdsaCurve == nil {
			ecdsaCurve = elliptic.P256()
		}
		return ecdsa.GenerateKey(ecdsaCurve, rand.Reader)
	case KeyTypeEd25519:
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		return priv, err
	default:
		return nil, ErrInvalidKeyPair
	}
}

func generateSerialNumber() *big.Int {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	return serialNumber
}

func GenerateCA(config *CAConfig) (*Certificate, error) {
	if config == nil {
		config = DefaultCAConfig()
	}

	privateKey, err := generatePrivateKey(config.KeyType, config.RSAKeySize, config.ECDSACurve)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	serialNumber := config.SerialNumber
	if serialNumber == nil {
		serialNumber = generateSerialNumber()
	}

	notBefore := config.NotBefore
	if notBefore.IsZero() {
		notBefore = time.Now()
	}

	notAfter := config.NotAfter
	if notAfter.IsZero() {
		notAfter = time.Now().AddDate(10, 0, 0)
	}

	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               config.Subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, publicKey(privateKey), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	keyPEM, err := encodePrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return &Certificate{
		Certificate: cert,
		PrivateKey:  privateKey,
		PublicKey:   publicKey(privateKey),
		CertPEM:     certPEM,
		KeyPEM:      keyPEM,
	}, nil
}

func GenerateServerCertificate(ca *Certificate, config *ServerCertConfig) (*Certificate, error) {
	if ca == nil {
		return nil, ErrInvalidCertificate
	}
	if config == nil {
		config = DefaultServerCertConfig()
	}

	privateKey, err := generatePrivateKey(config.KeyType, config.RSAKeySize, config.ECDSACurve)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	serialNumber := generateSerialNumber()
	notBefore := config.NotBefore
	if notBefore.IsZero() {
		notBefore = time.Now()
	}
	notAfter := config.NotAfter
	if notAfter.IsZero() {
		notAfter = time.Now().AddDate(1, 0, 0)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"BIOMETRICS CLI"},
			CommonName:   "BIOMETRICS CLI Server",
		},
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    config.DNSNames,
		IPAddresses: config.IPAddresses,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, ca.Certificate, publicKey(privateKey), ca.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	keyPEM, err := encodePrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return &Certificate{
		Certificate: cert,
		PrivateKey:  privateKey,
		PublicKey:   publicKey(privateKey),
		CertPEM:     certPEM,
		KeyPEM:      keyPEM,
	}, nil
}

func GenerateClientCertificate(ca *Certificate, config *ClientCertConfig) (*Certificate, error) {
	if ca == nil {
		return nil, ErrInvalidCertificate
	}
	if config == nil {
		config = DefaultClientCertConfig()
	}

	privateKey, err := generatePrivateKey(config.KeyType, config.RSAKeySize, config.ECDSACurve)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	serialNumber := generateSerialNumber()
	notBefore := config.NotBefore
	if notBefore.IsZero() {
		notBefore = time.Now()
	}
	notAfter := config.NotAfter
	if notAfter.IsZero() {
		notAfter = time.Now().AddDate(1, 0, 0)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      config.Subject,
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, ca.Certificate, publicKey(privateKey), ca.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	keyPEM, err := encodePrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return &Certificate{
		Certificate: cert,
		PrivateKey:  privateKey,
		PublicKey:   publicKey(privateKey),
		CertPEM:     certPEM,
		KeyPEM:      keyPEM,
	}, nil
}

func ValidateCertificate(cert *x509.Certificate, roots *x509.CertPool) error {
	if cert == nil {
		return ErrInvalidCertificate
	}

	now := time.Now()
	if now.Before(cert.NotBefore) {
		return ErrCertificateNotValid
	}
	if now.After(cert.NotAfter) {
		return ErrCertificateExpired
	}

	if roots != nil {
		opts := x509.VerifyOptions{
			Roots: roots,
		}
		if _, err := cert.Verify(opts); err != nil {
			return fmt.Errorf("certificate verification failed: %w", err)
		}
	}

	return nil
}

func VerifyKeyPair(cert *x509.Certificate, privateKey interface{}) error {
	if cert == nil || privateKey == nil {
		return ErrInvalidKeyPair
	}

	certPubKey := cert.PublicKey
	privPubKey := publicKey(privateKey)

	if !publicKeysEqual(certPubKey, privPubKey) {
		return ErrKeyMismatch
	}

	return nil
}

func publicKey(privateKey interface{}) interface{} {
	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func publicKeysEqual(a, b interface{}) bool {
	switch a := a.(type) {
	case *rsa.PublicKey:
		b, ok := b.(*rsa.PublicKey)
		if !ok {
			return false
		}
		return a.N.Cmp(b.N) == 0 && a.E == b.E
	case *ecdsa.PublicKey:
		b, ok := b.(*ecdsa.PublicKey)
		if !ok {
			return false
		}
		return a.Curve == b.Curve && a.X.Cmp(b.X) == 0 && a.Y.Cmp(b.Y) == 0
	case ed25519.PublicKey:
		b, ok := b.(ed25519.PublicKey)
		if !ok {
			return false
		}
		return a.Equal(b)
	default:
		return false
	}
}

func encodePrivateKey(privateKey interface{}) ([]byte, error) {
	var keyBytes []byte
	var keyType string

	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		keyBytes = x509.MarshalPKCS1PrivateKey(k)
		keyType = "RSA PRIVATE KEY"
	case *ecdsa.PrivateKey:
		var err error
		keyBytes, err = x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}
		keyType = "EC PRIVATE KEY"
	case ed25519.PrivateKey:
		var err error
		keyBytes, err = x509.MarshalPKCS8PrivateKey(k)
		if err != nil {
			return nil, err
		}
		keyType = "PRIVATE KEY"
	default:
		return nil, ErrInvalidKeyPair
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  keyType,
		Bytes: keyBytes,
	}), nil
}

func SaveCertificate(cert *Certificate, certPath, keyPath string) error {
	if cert == nil {
		return ErrInvalidCertificate
	}

	if err := os.MkdirAll(filepath.Dir(certPath), 0755); err != nil {
		return fmt.Errorf("failed to create certificate directory: %w", err)
	}

	if err := os.WriteFile(certPath, cert.CertPEM, 0644); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}

	if err := os.WriteFile(keyPath, cert.KeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

func LoadCertificate(certPath, keyPath string) (*Certificate, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, fmt.Errorf("failed to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to decode key PEM")
	}

	var privateKey interface{}
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		privateKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	case "EC PRIVATE KEY":
		privateKey, err = x509.ParseECPrivateKey(keyBlock.Bytes)
	case "PRIVATE KEY":
		privateKey, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	default:
		return nil, fmt.Errorf("unknown key type: %s", keyBlock.Type)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &Certificate{
		Certificate: cert,
		PrivateKey:  privateKey,
		PublicKey:   cert.PublicKey,
		CertPEM:     certPEM,
		KeyPEM:      keyPEM,
	}, nil
}

func CreateCertPool(certs ...*x509.Certificate) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, cert := range certs {
		pool.AddCert(cert)
	}
	return pool, nil
}
