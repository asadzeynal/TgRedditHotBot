package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

func Encrypt(value string, keyPhrase string) []byte {
	aesBlock, err := aes.NewCipher([]byte(keyPhrase))
	if err != nil {
		fmt.Println(err)
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	cipheredText := gcmInstance.Seal(nonce, nonce, []byte(value), nil)

	return cipheredText
}

func Decrypt(ciphered []byte, keyPhrase string) string {
	aesBlock, err := aes.NewCipher([]byte(keyPhrase))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(1)
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(2)
	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(3)
	return string(originalText)
}
