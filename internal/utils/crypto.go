package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/joho/godotenv"
)

// Переменные для ключей и секретов
var (
	pgpEntity  *openpgp.Entity
	hmacSecret []byte
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning loading .env: %v", err)
	}

	// Загружаем HMAC_SECRET из окружения
	hmacSecretStr := os.Getenv("HMAC_SECRET")
	if hmacSecretStr == "" {
		panic("HMAC_SECRET is not set in environment")
	}
	// if len(hmacSecretStr) < 32 {
	// 	log.Printf("Warning: HMAC_SECRET is shorter than recommended 32 bytes")
	// }
	hmacSecret = []byte(hmacSecretStr)

	// Загружаем PGP_PRIVATE_KEY из окружения
	privateKeyFile := os.Getenv("PGP_PRIVATE_KEY")
	privateKeyBytes, err := os.ReadFile(privateKeyFile)
	if err != nil {
		panic("failed to read PGP private key file")
	}
	privateKeyArmored := string(privateKeyBytes)
	if privateKeyArmored == "" {
		panic("PGP_PRIVATE_KEY is not set in environment")
	}

	entityList, err := readPGPEntityFromString(privateKeyArmored)
	if err != nil {
		panic("failed to parse PGP private key")
	}
	if len(entityList) == 0 {
		panic("no PGP entities found in key")
	}
	pgpEntity = entityList[0]
	
}

func readPGPEntityFromString(key string) ([]*openpgp.Entity, error) {
	block, err := armor.Decode(strings.NewReader(key))
	if err != nil {
		return nil, fmt.Errorf("failed to decode armored key: %w", err)
	}
	entity, err := openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to read entity: %w", err)
	}
	return []*openpgp.Entity{entity}, nil
}

func EncryptPGP(data string) (string, error) {
	var buf bytes.Buffer
	w, err := armor.Encode(&buf, "PGP MESSAGE", nil)
	if err != nil {
		return "", err
	}
	plaintext, err := openpgp.Encrypt(w, []*openpgp.Entity{pgpEntity}, nil, nil, nil)
	if err != nil {
		return "", err
	}
	_, err = plaintext.Write([]byte(data))
	if err != nil {
		return "", err
	}
	plaintext.Close()
	w.Close()
	return buf.String(), nil
}

func DecryptPGP(encrypted string) (string, error) {
	block, err := armor.Decode(strings.NewReader(encrypted))
	if err != nil {
		return "", err
	}
	md, err := openpgp.ReadMessage(block.Body, openpgp.EntityList{pgpEntity}, nil, nil)
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ComputeHMAC(data string) string {
	h := hmac.New(sha256.New, hmacSecret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyHMAC(data, mac []byte) bool {
	expectedMAC := ComputeHMAC(string(data))
	return hmac.Equal([]byte(expectedMAC), mac)
}
