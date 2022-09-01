package toolkit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var keyText = "Person is a thinking reed!"
var nonceText = "Recognize yourself!"

// Encrypt by md5
func EncryptByMD5(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	encrypted := fmt.Sprintf("%x", h.Sum(nil))

	return encrypted
}

// Generate key
func GenerateAESKey() string {
	encryptedKey := EncryptByMD5(keyText)
	encryptedNonce := EncryptByMD5(nonceText)

	return encryptedKey + encryptedNonce
}

// Generate nonce
func GenerateAESNonce() string {
	encrypted := EncryptByMD5(nonceText)
	subEncrypted := encrypted[:24]

	return subEncrypted
}

// Encrypt
func Encrypt(text string) string {
	key, _ := hex.DecodeString(GenerateAESKey())
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	nonce, _ := hex.DecodeString(GenerateAESNonce())

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	formatCiphertext := fmt.Sprintf("%x", ciphertext)

	return formatCiphertext
}

// Decrypt
func Decrypt(text string) string {
	key, _ := hex.DecodeString(GenerateAESKey())
	ciphertext, _ := hex.DecodeString(text)
	nonce, _ := hex.DecodeString(GenerateAESNonce())

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	formatPlaintext := string(plaintext)

	return formatPlaintext
}

// Get token
func GetToken(r *http.Request) (string, string) {
	vars := r.URL.Query()
	token := vars["token"]
	splitedToken := strings.Split(token[0], "@")

	id := Decrypt(splitedToken[0])
	timestamp := Decrypt(splitedToken[1])

	return id, timestamp
}
