package crypro

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"hash"
	"io"
	"os"
)

const (
	privateKeyPEMBlockType = "RSA PRIVATE KEY"
	publicKeyPEMBlockType  = "RSA PUBLIC KEY"
)

type CManager struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func NewCManager() *CManager {
	return &CManager{}
}

func (cm *CManager) InitPublicKey(pathCryptoKey string) error {
	if pathCryptoKey != "" {
		key, err := getPublicKey(pathCryptoKey)
		if err != nil {
			return fmt.Errorf("failed to read PEM file: %v", err)
		}
		cm.PublicKey = key
	}

	return nil
}

func (cm *CManager) InitPrivateKey(pathCryptoKey string) error {
	if pathCryptoKey != "" {
		key, err := getPrivateKey(pathCryptoKey)
		if err != nil {
			return fmt.Errorf("failed to read PEM file: %v", err)
		}
		cm.PrivateKey = key
	}

	return nil
}

func getPrivateKey(pathCryptoKey string) (*rsa.PrivateKey, error) {
	pemBlock, err := getPemBlock(pathCryptoKey)
	if err != nil {
		return nil, err
	}
	if pemBlock.Type != privateKeyPEMBlockType {
		return nil, fmt.Errorf("unsupported PEM block type: %s", pemBlock.Type)
	}

	return x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
}

func getPublicKey(pathCryptoKey string) (*rsa.PublicKey, error) {
	pemBlock, err := getPemBlock(pathCryptoKey)
	if err != nil {
		return nil, err
	}
	if pemBlock.Type != publicKeyPEMBlockType {
		return nil, fmt.Errorf("unsupported PEM block type: %s", pemBlock.Type)
	}

	return x509.ParsePKCS1PublicKey(pemBlock.Bytes)
}

func getPemBlock(pathCryptoKey string) (*pem.Block, error) {
	keyPEM, err := os.ReadFile(pathCryptoKey)
	if err != nil {
		return nil, fmt.Errorf("failed cryptoKey ReadFile: %v", err)
	}
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, fmt.Errorf("failed pem.Decode")
	}

	return keyBlock, nil
}

func (cm *CManager) Decrypt(bodyBytes []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("cannot decode base64: %v", err)
	}

	encryptedBytes, err := decryptOAEP(sha256.New(), rand.Reader, cm.PrivateKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot decryptOAEP: %v", err)
	}

	return encryptedBytes, nil
}

func (cm *CManager) Encrypt(body *bytes.Buffer) (*bytes.Buffer, error) {
	if cm.PublicKey == nil {
		return body, nil
	}

	encryptedBytes, err := encryptOAEP(sha256.New(), rand.Reader, cm.PublicKey, body.Bytes(), nil)
	if err != nil {
		return nil, fmt.Errorf("rsa.EncryptOAEP: %v", err)
	}

	return bytes.NewBuffer([]byte(base64.StdEncoding.EncodeToString(encryptedBytes))), nil
}

func encryptOAEP(hash hash.Hash, random io.Reader, public *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func decryptOAEP(hash hash.Hash, random io.Reader, private *rsa.PrivateKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, private, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
