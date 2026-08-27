package idcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	KeyBytes   = 32
	version    = 0x01
	nonceBytes = 12
)

var ErrCorrupt = errors.New("idcrypto: ciphertext เสียหายหรือ key ไม่ตรง")

type Cipher struct {
	aead cipher.AEAD
}

func New(hexKey string) (Cipher, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return Cipher{}, fmt.Errorf("idcrypto: key ไม่ใช่ hex: %w", err)
	}
	if len(key) != KeyBytes {
		return Cipher{}, fmt.Errorf("idcrypto: key ต้องยาว %d ไบต์ (hex %d ตัว) ได้ %d", KeyBytes, KeyBytes*2, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return Cipher{}, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return Cipher{}, err
	}
	return Cipher{aead: aead}, nil
}

func (c Cipher) Seal(plaintext, aad string) ([]byte, error) {
	nonce := make([]byte, nonceBytes)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	out := make([]byte, 0, 1+nonceBytes+len(plaintext)+c.aead.Overhead())
	out = append(out, version)
	out = append(out, nonce...)
	return c.aead.Seal(out, nonce, []byte(plaintext), []byte(aad)), nil
}

func (c Cipher) Open(b []byte, aad string) (string, error) {
	if len(b) < 1+nonceBytes || b[0] != version {
		return "", ErrCorrupt
	}
	plain, err := c.aead.Open(nil, b[1:1+nonceBytes], b[1+nonceBytes:], []byte(aad))
	if err != nil {
		return "", ErrCorrupt
	}
	return string(plain), nil
}
