package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

var key = []byte{6, 189, 106, 125, 221, 172, 17, 103, 153, 126, 87, 44, 31, 169, 153, 64,
	133, 62, 137, 100, 236, 28, 198, 20, 153, 191, 214, 111, 146, 138, 144, 126}

func Encrypt(b []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, 
			fmt.Errorf("failed to create cipher block: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, 
			fmt.Errorf("failed to create aesgcm: %w", err)
	}

	res := aesgcm.Seal(b[:0], make([]byte, aesgcm.NonceSize()), b, nil)

	return res, nil
}

func Decrypt(b []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, 
			fmt.Errorf("failed to create cipher block: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, 
			fmt.Errorf("failed to create aesgcm: %w", err)
	}

	res, err := aesgcm.Open(b[:0], make([]byte, aesgcm.NonceSize()), b, nil)
	if err != nil {
		return nil, 
			fmt.Errorf("failed to create aesgcm: %w", err)
	}

	return res, nil
}