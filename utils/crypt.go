package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"

	"golang.org/x/crypto/bcrypt"
)

// Initiates a DES key, pass a string to get an 8 byte key.
func InitDESKey(key string) [8]byte {
	result := make([]byte, 8)
	for i, j := 0, 0; i < 40; i++ {
		if j < len(key) {
			result[i%8] ^= key[j]
		} else {
			result[i%8] ^= 0x00
		}
		j++
	}
	return [8]byte(result)
}

// Encrypts a string based on bcrypt algorithm.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Verifies if a hash matches the bcrypt password.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Verifies if a hash matches the md5 password.
func VerifyPasswordMD5(password, hash string) bool {
	md5Bytes := md5.Sum([]byte(password))
	return hex.EncodeToString(md5Bytes[:]) == hash
}

func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	b := block.Bytes
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}
	return key, nil
}

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func PKCS5Trimming(encrypt []byte) []byte {
	if len(encrypt) == 0 {
		return encrypt
	}
	padding := encrypt[len(encrypt)-1]
	if len(encrypt)-int(padding) < 0 {
		return encrypt
	}
	return encrypt[:len(encrypt)-int(padding)]
}
