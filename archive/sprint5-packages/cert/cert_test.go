package cert

import (
	"crypto/elliptic"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGenerateCA(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		cert, err := GenerateCA(nil)
		if err != nil {
			t.Fatalf("GenerateCA() error = %v", err)
		}
		if cert == nil {
			t.Fatal("GenerateCA() returned nil certificate")
		}
		if cert.Certificate == nil {
			t.Error("Certificate is nil")
		}
		if cert.PrivateKey == nil {
			t.Error("PrivateKey is nil")
		}
		if cert.PublicKey == nil {
			t.Error("PublicKey is nil")
		}
		if len(cert.CertPEM) == 0 {
			t.Error("CertPEM is empty")
		}
		if len(cert.KeyPEM) == 0 {
			t.Error("KeyPEM is empty")
		}
		if !cert.Certificate.IsCA {
			t.Error("Certificate should be a CA")
		}
	})

	t.Run("CustomConfig", func(t *testing.T) {
		config := &CAConfig{
			Subject: pkix.Name{
				Organization: []string{"Test Org"},
				CommonName:   "Test CA",
			},
			KeyType:    KeyTypeRSA,
			RSAKeySize: 2048,
			NotBefore:  time.Now(),
			NotAfter:   time.Now().AddDate(5, 0, 0),
		}

		cert, err := GenerateCA(config)
		if err != nil {
			t.Fatalf("GenerateCA() error = %v", err)
		}

		if cert.Certificate.Subject.CommonName != "Test CA" {
			t.Errorf("CommonName = %v, want Test CA", cert.Certificate.Subject.CommonName)
		}
	})

	t.Run("ECDSAKey", func(t *testing.T) {
		config := &CAConfig{
			Subject:    pkix.Name{CommonName: "ECDSA CA"},
			KeyType:    KeyTypeECDSA,
			ECDSACurve: elliptic.P384(),
			NotBefore:  time.Now(),
			NotAfter:   time.Now().AddDate(1, 0, 0),
		}

		cert, err := GenerateCA(config)
		if err != nil {
			t.Fatalf("GenerateCA() error = %v", err)
		}

		if cert.Certificate.PublicKeyAlgorithm != x509.ECDSA {
			t.Errorf("PublicKeyAlgorithm = %v, want ECDSA", cert.Certificate.PublicKeyAlgorithm)
		}
	})

	t.Run("Ed25519Key", func(t *testing.T) {
		config := &CAConfig{
			Subject:   pkix.Name{CommonName: "Ed25519 CA"},
			KeyType:   KeyTypeEd25519,
			NotBefore: time.Now(),
			NotAfter:  time.Now().AddDate(1, 0, 0),
		}

		cert, err := GenerateCA(config)
		if err != nil {
			t.Fatalf("GenerateCA() error = %v", err)
		}

		if cert.Certificate.PublicKeyAlgorithm != x509.Ed25519 {
			t.Errorf("PublicKeyAlgorithm = %v, want Ed25519", cert.Certificate.PublicKeyAlgorithm)
		}
	})
}

func TestGenerateServerCertificate(t *testing.T) {
	ca, err := GenerateCA(nil)
	if err != nil {
		t.Fatalf("Failed to generate CA: %v", err)
	}

	t.Run("DefaultConfig", func(t *testing.T) {
		cert, err := GenerateServerCertificate(ca, nil)
		if err != nil {
			t.Fatalf("GenerateServerCertificate() error = %v", err)
		}
		if cert == nil {
			t.Fatal("GenerateServerCertificate() returned nil")
		}
		if cert.Certificate.IsCA {
			t.Error("Server certificate should not be a CA")
		}
	})

	t.Run("WithDNSNames", func(t *testing.T) {
		config := &ServerCertConfig{
			DNSNames:    []string{"example.com", "api.example.com"},
			IPAddresses: []net.IP{net.ParseIP("192.168.1.1")},
			NotBefore:   time.Now(),
			NotAfter:    time.Now().AddDate(1, 0, 0),
		}

		cert, err := GenerateServerCertificate(ca, config)
		if err != nil {
			t.Fatalf("GenerateServerCertificate() error = %v", err)
		}

		if len(cert.Certificate.DNSNames) != 2 {
			t.Errorf("DNSNames count = %d, want 2", len(cert.Certificate.DNSNames))
		}
		if len(cert.Certificate.IPAddresses) != 1 {
			t.Errorf("IPAddresses count = %d, want 1", len(cert.Certificate.IPAddresses))
		}
	})

	t.Run("NilCA", func(t *testing.T) {
		_, err := GenerateServerCertificate(nil, nil)
		if err == nil {
			t.Error("GenerateServerCertificate() should fail with nil CA")
		}
	})
}

func TestGenerateClientCertificate(t *testing.T) {
	ca, err := GenerateCA(nil)
	if err != nil {
		t.Fatalf("Failed to generate CA: %v", err)
	}

	t.Run("DefaultConfig", func(t *testing.T) {
		cert, err := GenerateClientCertificate(ca, nil)
		if err != nil {
			t.Fatalf("GenerateClientCertificate() error = %v", err)
		}
		if cert == nil {
			t.Fatal("GenerateClientCertificate() returned nil")
		}
		if cert.Certificate.IsCA {
			t.Error("Client certificate should not be a CA")
		}
	})

	t.Run("WithSubject", func(t *testing.T) {
		config := &ClientCertConfig{
			Subject: pkix.Name{
				Organization: []string{"Test Client Org"},
				CommonName:   "test-client",
			},
			NotBefore: time.Now(),
			NotAfter:  time.Now().AddDate(1, 0, 0),
		}

		cert, err := GenerateClientCertificate(ca, config)
		if err != nil {
			t.Fatalf("GenerateClientCertificate() error = %v", err)
		}

		if cert.Certificate.Subject.CommonName != "test-client" {
			t.Errorf("CommonName = %v, want test-client", cert.Certificate.Subject.CommonName)
		}
	})

	t.Run("NilCA", func(t *testing.T) {
		_, err := GenerateClientCertificate(nil, nil)
		if err == nil {
			t.Error("GenerateClientCertificate() should fail with nil CA")
		}
	})
}

func TestValidateCertificate(t *testing.T) {
	ca, _ := GenerateCA(nil)

	t.Run("ValidCertificate", func(t *testing.T) {
		err := ValidateCertificate(ca.Certificate, nil)
		if err != nil {
			t.Errorf("ValidateCertificate() error = %v", err)
		}
	})

	t.Run("ExpiredCertificate", func(t *testing.T) {
		config := &CAConfig{
			Subject:   pkix.Name{CommonName: "Expired CA"},
			NotBefore: time.Now().Add(-24 * time.Hour),
			NotAfter:  time.Now().Add(-1 * time.Hour),
		}
		cert, _ := GenerateCA(config)

		err := ValidateCertificate(cert.Certificate, nil)
		if err != ErrCertificateExpired {
			t.Errorf("ValidateCertificate() error = %v, want ErrCertificateExpired", err)
		}
	})

	t.Run("NotYetValid", func(t *testing.T) {
		config := &CAConfig{
			Subject:   pkix.Name{CommonName: "Future CA"},
			NotBefore: time.Now().Add(24 * time.Hour),
			NotAfter:  time.Now().AddDate(1, 0, 0),
		}
		cert, _ := GenerateCA(config)

		err := ValidateCertificate(cert.Certificate, nil)
		if err != ErrCertificateNotValid {
			t.Errorf("ValidateCertificate() error = %v, want ErrCertificateNotValid", err)
		}
	})

	t.Run("NilCertificate", func(t *testing.T) {
		err := ValidateCertificate(nil, nil)
		if err != ErrInvalidCertificate {
			t.Errorf("ValidateCertificate() error = %v, want ErrInvalidCertificate", err)
		}
	})
}

func TestVerifyKeyPair(t *testing.T) {
	ca, _ := GenerateCA(nil)

	t.Run("ValidPair", func(t *testing.T) {
		err := VerifyKeyPair(ca.Certificate, ca.PrivateKey)
		if err != nil {
			t.Errorf("VerifyKeyPair() error = %v", err)
		}
	})

	t.Run("MismatchedPair", func(t *testing.T) {
		otherCA, _ := GenerateCA(nil)
		err := VerifyKeyPair(ca.Certificate, otherCA.PrivateKey)
		if err != ErrKeyMismatch {
			t.Errorf("VerifyKeyPair() error = %v, want ErrKeyMismatch", err)
		}
	})

	t.Run("NilCertificate", func(t *testing.T) {
		err := VerifyKeyPair(nil, ca.PrivateKey)
		if err != ErrInvalidKeyPair {
			t.Errorf("VerifyKeyPair() error = %v, want ErrInvalidKeyPair", err)
		}
	})

	t.Run("NilPrivateKey", func(t *testing.T) {
		err := VerifyKeyPair(ca.Certificate, nil)
		if err != ErrInvalidKeyPair {
			t.Errorf("VerifyKeyPair() error = %v, want ErrInvalidKeyPair", err)
		}
	})
}

func TestSaveAndLoadCertificate(t *testing.T) {
	ca, _ := GenerateCA(nil)

	tmpDir, err := os.MkdirTemp("", "cert-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	certPath := filepath.Join(tmpDir, "ca.crt")
	keyPath := filepath.Join(tmpDir, "ca.key")

	t.Run("SaveAndLoad", func(t *testing.T) {
		err := SaveCertificate(ca, certPath, keyPath)
		if err != nil {
			t.Fatalf("SaveCertificate() error = %v", err)
		}

		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			t.Error("Certificate file was not created")
		}
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			t.Error("Key file was not created")
		}

		loaded, err := LoadCertificate(certPath, keyPath)
		if err != nil {
			t.Fatalf("LoadCertificate() error = %v", err)
		}

		if loaded.Certificate.Subject.CommonName != ca.Certificate.Subject.CommonName {
			t.Errorf("Loaded CommonName = %v, want %v",
				loaded.Certificate.Subject.CommonName,
				ca.Certificate.Subject.CommonName)
		}
	})

	t.Run("NilCertificate", func(t *testing.T) {
		err := SaveCertificate(nil, certPath, keyPath)
		if err == nil {
			t.Error("SaveCertificate() should fail with nil certificate")
		}
	})
}

func TestCreateCertPool(t *testing.T) {
	ca1, _ := GenerateCA(&CAConfig{Subject: pkix.Name{CommonName: "CA1"}})
	ca2, _ := GenerateCA(&CAConfig{Subject: pkix.Name{CommonName: "CA2"}})

	pool, err := CreateCertPool(ca1.Certificate, ca2.Certificate)
	if err != nil {
		t.Fatalf("CreateCertPool() error = %v", err)
	}
	if pool == nil {
		t.Fatal("CreateCertPool() returned nil")
	}
}

func TestDefaultConfigs(t *testing.T) {
	t.Run("DefaultCAConfig", func(t *testing.T) {
		config := DefaultCAConfig()
		if config == nil {
			t.Fatal("DefaultCAConfig() returned nil")
		}
		if config.KeyType != KeyTypeECDSA {
			t.Errorf("KeyType = %v, want KeyTypeECDSA", config.KeyType)
		}
	})

	t.Run("DefaultServerCertConfig", func(t *testing.T) {
		config := DefaultServerCertConfig()
		if config == nil {
			t.Fatal("DefaultServerCertConfig() returned nil")
		}
		if len(config.DNSNames) == 0 {
			t.Error("DNSNames should not be empty")
		}
	})

	t.Run("DefaultClientCertConfig", func(t *testing.T) {
		config := DefaultClientCertConfig()
		if config == nil {
			t.Fatal("DefaultClientCertConfig() returned nil")
		}
		if config.Subject.CommonName == "" {
			t.Error("CommonName should not be empty")
		}
	})
}
