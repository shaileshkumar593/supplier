package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

// AESEncrypt - Encrypt string using AES with block size of 16
func AESEncrypt(text, secretKey string) string {
	var (
		key []byte
	)

	key, _ = hex.DecodeString(secretKey)

	plaintext := []byte(text)

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	t := aesgcm.Seal(nil, nonce, plaintext, nil)

	// pass the nonce for decryption. Its the last 12 bytes of the hex string
	return fmt.Sprintf("%x%x", t, nonce)
}

// AESDecrypt - Decrypt string using AES with block size of 16
func AESDecrypt(ciphertext, secretKey string) string {

	var (
		key, cbytes, nonce []byte
		block              cipher.Block
		err                error
	)

	key, _ = hex.DecodeString(secretKey)

	// trim last 24 character to get the cipher text
	cbytes, _ = hex.DecodeString(ciphertext[0 : len(ciphertext)-24])
	// get the last 24 character which is the nonce
	nonce, _ = hex.DecodeString(ciphertext[len(ciphertext)-24:])
	block, err = aes.NewCipher(key)

	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, cbytes, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

// AESDecryptWithIV decrypts provided encrypted AES 256 CBC data
// and IV is parsed from the data received
func AESDecryptWithIV(data, secretKey string) (string, error) {
	encryptedData := strings.Split(data, ":")
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData[0])
	if err != nil {
		return "", err
	}

	key := []byte(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short")
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, []byte(encryptedData[1]))
	mode.CryptBlocks(cipherText, cipherText)

	return string(PKCS5Trimming(cipherText)), nil
}

// PKCS5Trimming creating own padding scheme to trim unnecessary values
func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
