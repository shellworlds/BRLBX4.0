package devcert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// SignCSR issues a client-auth TLS certificate from a PEM CA cert + RSA private key.
func SignCSR(csrPEM, caCertPEM, caKeyPEM []byte) (certPEM []byte, serialHex string, err error) {
	csrBlock, _ := pem.Decode(csrPEM)
	if csrBlock == nil {
		return nil, "", fmt.Errorf("invalid csr pem")
	}
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return nil, "", err
	}
	if err := csr.CheckSignature(); err != nil {
		return nil, "", fmt.Errorf("csr signature: %w", err)
	}

	caBlock, _ := pem.Decode(caCertPEM)
	if caBlock == nil {
		return nil, "", fmt.Errorf("invalid ca cert pem")
	}
	ca, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, "", err
	}

	keyBlock, _ := pem.Decode(caKeyPEM)
	if keyBlock == nil {
		return nil, "", fmt.Errorf("invalid ca key pem")
	}
	var caKey *rsa.PrivateKey
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		k, e := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if e != nil {
			return nil, "", e
		}
		caKey = k
	case "PRIVATE KEY":
		k, e := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
		if e != nil {
			return nil, "", e
		}
		var ok bool
		caKey, ok = k.(*rsa.PrivateKey)
		if !ok {
			return nil, "", fmt.Errorf("ca key must be RSA for this signer")
		}
	default:
		return nil, "", fmt.Errorf("unsupported pem type %q", keyBlock.Type)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, "", err
	}
	serialHex = fmt.Sprintf("%x", serial)

	tpl := x509.Certificate{
		SerialNumber:   serial,
		Subject:        csr.Subject,
		NotBefore:      time.Now().Add(-1 * time.Hour),
		NotAfter:       time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:       x509.KeyUsageDigitalSignature,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		DNSNames:       csr.DNSNames,
		IPAddresses:    csr.IPAddresses,
		EmailAddresses: csr.EmailAddresses,
		URIs:           csr.URIs,
	}

	der, err := x509.CreateCertificate(rand.Reader, &tpl, ca, csr.PublicKey, caKey)
	if err != nil {
		return nil, "", err
	}
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	return certPEM, serialHex, nil
}
