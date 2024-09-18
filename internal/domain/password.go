package domain

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

// PasswordVerifier compares a password with a hashed password.
type PasswordVerifier interface {
	Verify(password string, hashed HashedPassword) (bool, error)
}

type PasswordVerifierFunc func(password string, hashed HashedPassword) (bool, error)

func (f PasswordVerifierFunc) Verify(password string, hashed HashedPassword) (bool, error) {
	return f(password, hashed)
}

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// Based on https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
var defaultParams = &params{
	memory:      64 * 1024,
	iterations:  1,
	parallelism: 1,
	saltLength:  16,
	keyLength:   32,
}

// Argon2idHashPassword hashes a password using Argon2id.
func Argon2idHashPassword(password string) (HashedPassword, error) {
	// Generate a random salt.
	salt := make([]byte, defaultParams.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the password.
	hash := argon2.IDKey([]byte(password), salt, defaultParams.iterations, defaultParams.memory, defaultParams.parallelism, defaultParams.keyLength)

	// Encode the parameters, salt, and hash into a string.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, defaultParams.memory, defaultParams.iterations, defaultParams.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash))

	return HashedPassword(encodedHash), nil
}

// Argon2idVerifyPassword implements PasswordVerifier.
var Argon2idVerifyPassword PasswordVerifierFunc = argon2idVerifyPassword

// argon2idVerifyPassword checks if a password matches a hash.
func argon2idVerifyPassword(password string, encodedHash HashedPassword) (bool, error) {
	// Extract the parameters, salt, and hash from the encoded string.
	parts := strings.Split(string(encodedHash), "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var p params
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return false, fmt.Errorf("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("invalid salt")
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("invalid hash")
	}

	p.keyLength = uint32(len(hash))

	// Compute the hash of the provided password.
	computedHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Compare the computed hash with the stored hash.
	return string(hash) == string(computedHash), nil
}
