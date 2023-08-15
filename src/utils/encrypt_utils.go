package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func EncryptString(s string) string {
	firstEncryptData := sha256.Sum256([]byte(s))
	firstEncryptString := fmt.Sprintf("%x", firstEncryptData)
	lastEncryptData := sha256.Sum256([]byte(fmt.Sprintf("%s%s%s", os.Getenv("PRE_SALT"), firstEncryptString, os.Getenv("END_SALT"))))

	return fmt.Sprintf("%x", lastEncryptData)
}

func GeneratePassword(complexities, passwordLength int) string {
	const (
		lowerLetter      = "abcdefghijklmnopqrstuvwxyz"
		upperLetter      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digit            = "0123456789"
		specialCharacter = "!@#$%^&*?"
	)
	var conditions []string
	if complexities&8 > 0 {
		conditions = append(conditions, lowerLetter)
	}
	if complexities&4 > 0 {
		conditions = append(conditions, upperLetter)
	}
	if complexities&2 > 0 {
		conditions = append(conditions, digit)
	}
	if complexities&1 > 0 {
		conditions = append(conditions, specialCharacter)
	}

	rand.Seed(ConvertToNanoTime(GetCurrentTime()))

	password := make([]byte, passwordLength)
	for _, charset := range conditions {
		index := rand.Intn(passwordLength)
		for password[index] != 0 {
			index = rand.Intn(passwordLength)
		}
		char := charset[rand.Intn(len(charset))]
		password[index] = char
	}

	for i := 0; i < passwordLength; i++ {
		if password[i] == 0 {
			charset := conditions[rand.Intn(len(conditions))]
			password[i] = charset[rand.Intn(len(charset))]
		}
	}

	return string(password)
}

func EncryptPlainText(plainText []byte, salt string) string {
	key := generateKey(salt)
	iv := generateIV(salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	plainText = PKCS7Padding(plainText, block.BlockSize())

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(plainText))

	return Base64EncodeString(ciphertext)
}

func DecryptCipherText(cipherText, salt string) string {
	key := generateKey(salt)
	iv := generateIV(salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertextBytes := Base64DecodeString(cipherText)

	plainText := make([]byte, len(ciphertextBytes)-aes.BlockSize)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, ciphertextBytes[aes.BlockSize:])

	plainText = PKCS7Unpadding(plainText)

	return string(plainText)
}

func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := strings.Repeat(string(byte(padding)), padding)
	return append(data, []byte(padText)...)
}

func PKCS7Unpadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

func generateKey(salt string) []byte {
	key := make([]byte, 32)
	copy(key[:16], salt)
	copy(key[16:], salt)
	return key
}

func generateIV(salt string) []byte {
	iv := make([]byte, aes.BlockSize)
	copy(iv, salt)
	return iv
}
