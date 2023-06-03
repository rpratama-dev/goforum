package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Specify the custom file paths
var privateKeyPath = filepath.Join("./src/configs", "key-pair", "private.key")
var publicKeyPath = filepath.Join("./src/configs", "key-pair", "public.crt")

func isFileExist(pathToFile string) (bool)  {
	_, err := os.Stat(pathToFile)

	if err == nil {
  // path/to/whatever exists
		return true
	} 
	if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		return false
	} 
	// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
	// Schrodinger: file may or may not exist. See err for details.
	return false
}

// Generate public key and private key
func GenerateKeyPair() (bool) {
	isPrivateExist := isFileExist(privateKeyPath)
	isPublicExist := isFileExist(privateKeyPath)

	if (isPrivateExist && isPublicExist) {
		return true
	}

	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Failed to generate private key:", err)
	}

	// Save private key to file
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		log.Fatal("Failed to create private key file:", err)
	}
	defer privateKeyFile.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		log.Fatal("Failed to encode private key:", err)
	}

	fmt.Println("Private key saved to private.key")

	// Extract public key from private key
	publicKey := &privateKey.PublicKey

	// Save public key to file
	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		log.Fatal("Failed to create public key file:", err)
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal("Failed to marshal public key:", err)
	}

	publicKeyPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: publicKeyBytes,
	}

	err = pem.Encode(publicKeyFile, publicKeyPEM)
	if err != nil {
		log.Fatal("Failed to encode public key:", err)
	}

	fmt.Println("Public key saved to public.crt")
	return true
}

func GetPublicKey() (string)  {
	publicKey := ReadFile(publicKeyPath)
	if (publicKey == nil) {
		return ""
	}
	return string(publicKey)
}

func GetPrivateKey() (string)  {
	privateKey := ReadFile(privateKeyPath)
	if (privateKey == nil) {
		return ""
	}
	return string(privateKey)
}
