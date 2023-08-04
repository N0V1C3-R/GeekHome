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

func GeneratePassword() string {
	const (
		passwordLength   = 12
		lowerLetter      = "abcdefghijklmnopqrstuvwxyz"
		upperLetter      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digit            = "0123456789"
		specialCharacter = "!@#$%^&*_+<>?,.:"
	)
	rand.Seed(ConvertToNanoTime(GetCurrentTime()))

	password := make([]byte, passwordLength)
	password[rand.Intn(passwordLength)] = lowerLetter[rand.Intn(len(lowerLetter))]
	password[rand.Intn(passwordLength)] = upperLetter[rand.Intn(len(upperLetter))]
	password[rand.Intn(passwordLength)] = digit[rand.Intn(len(digit))]
	password[rand.Intn(passwordLength)] = specialCharacter[rand.Intn(len(specialCharacter))]

	for i := 0; i < passwordLength; i++ {
		if password[i] == 0 {
			charset := lowerLetter + upperLetter + digit + specialCharacter
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
	copy(key[:16], []byte(salt))
	copy(key[16:], []byte(salt))
	return key
}

func generateIV(salt string) []byte {
	iv := make([]byte, aes.BlockSize)
	copy(iv, []byte(salt))
	return iv
}
