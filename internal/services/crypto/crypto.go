package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log/slog"
)

// CryptoService сервис для шифрования строковых данных
type CryptoService interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(cipherText string) ([]byte, error)
}

type CryptoServiceImpl struct {
	log    *slog.Logger
	key    []byte
	cipher cipher.Block
}

func NewCryptoService(log *slog.Logger, key []byte) *CryptoServiceImpl {
	c, err := aes.NewCipher(key[:16])
	if err != nil {
		panic(err)
	}
	return &CryptoServiceImpl{
		log:    log,
		key:    key,
		cipher: c,
	}
}

// Encrypt метод для шифрования строки
func (s *CryptoServiceImpl) Encrypt(plaintext []byte) (string, error) {
	tag := "Encrypt"
	log := s.log.With("tag", tag)
	b := plaintext
	b = addPadding(b, aes.BlockSize)
	encMessage := make([]byte, len(b))
	iv := s.key[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(s.cipher, iv)
	mode.CryptBlocks(encMessage, b)

	encMessageString := base64.StdEncoding.EncodeToString(encMessage)

	log.Info("encrypted message" + encMessageString)
	return encMessageString, nil
}

// Decrypt метод для дешифрования строки
func (s *CryptoServiceImpl) Decrypt(cipherText string) ([]byte, error) {
	tag := "Decrypt"
	log := s.log.With("tag", tag)
	iv := s.key[:aes.BlockSize]

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		log.Error("base64 decoding error" + err.Error())
		return nil, err
	}

	if len(cipherTextBytes) < aes.BlockSize {
		log.Error("cipherText too short")
		return nil, errors.New("encMessage слишком короткий")
	}

	decrypted := make([]byte, len(cipherTextBytes))
	mode := cipher.NewCBCDecrypter(s.cipher, iv)

	mode.CryptBlocks(decrypted, cipherTextBytes)

	return removePadding(decrypted), nil
}

func addPadding(cipher []byte, blockSize int) []byte {
	padding := blockSize - len(cipher)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(cipher, padText...)
}

func removePadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])

	return src[:(length - unPadding)]
}
