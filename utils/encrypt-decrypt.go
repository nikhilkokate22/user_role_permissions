package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

var aesKey = []byte("12345678901234567890123456789012")
var fixedIV = []byte("1234567890123456")

// internal key loader
// func getAESKey() ([]byte, error) {
// 	secret := os.Getenv("SECRET_KEY")
// 	if secret == "" {
// 		return nil, errors.New("SECRET_KEY not set")
// 	}

// 	key := []byte(secret)
// 	if len(key) < 32 {
// 		return nil, errors.New("SECRET_KEY must be at least 32 bytes")
// 	}

// 	return key[:32], nil
// }

func EncryptAES(text string) (string, error) {

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()

	// PKCS7 padding
	padding := blockSize - len(text)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	plainText := append([]byte(text), padText...)

	mode := cipher.NewCBCEncrypter(block, fixedIV)

	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptAES(enc string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, fixedIV)

	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	// remove padding
	padding := int(plainText[len(plainText)-1])
	plainText = plainText[:len(plainText)-padding]

	return string(plainText), nil
}

// HashPassword hashes password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword(
// 		[]byte(password),
// 		bcrypt.DefaultCost,
// 	)
// 	return string(bytes), err
// }

// CheckPassword compares hash with plain password
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
