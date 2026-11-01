package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	memory  uint32
	time    uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func NewHashConfig() *HashConfig {
	return &HashConfig{
		memory:  64 * 1024,
		time:    2,
		threads: 2,
		keyLen:  32,
		saltLen: 16,
	}
}

func (h *HashConfig) GenerateSalt() []byte {
	salt := make([]byte, h.saltLen)

	rand.Read(salt)
	return salt
}

func (h *HashConfig) GenerateHash(password string) string {
	salt := h.GenerateSalt()
	hash := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)

	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	completeHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.memory, h.time, h.threads, base64Salt, base64Hash)

	return completeHash
}

func CompareHash(password, hashedPassword string) error {
	var version int
	var memory, time uint32
	var threads uint8

	result := strings.Split(hashedPassword, "$")
	if len(result) != 6 {
		return app_error.IncorrectHashFormat
	}
	if result[1] != "argon2id" {
		return app_error.IncorrectHashType
	}
	if _, err := fmt.Sscanf(result[2], "v=%d", &version); err != nil {
		return err
	}
	if version != argon2.Version {
		return app_error.IncorrectHashVersion
	}
	if _, err := fmt.Sscanf(result[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return err
	}

	salt, err := base64.RawStdEncoding.DecodeString(result[4])
	if err != nil {
		return err
	}

	hash, err := base64.RawStdEncoding.DecodeString(result[5])
	if err != nil {
		return err
	}

	newHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))
	if subtle.ConstantTimeCompare(hash, newHash) == 0 {
		return app_error.MismatchHash
	}

	return nil
}
