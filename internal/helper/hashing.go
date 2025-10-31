package helper

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memory     = 64 * 1024 // 64 MB
	iterations = 3
	threads    = 1
	saltLength = 16
	hashLength = 32
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, threads, hashLength)

	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$kl=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		threads,
		hashLength,
		b64salt,
		b64hash,
	), nil
}

func VerifyPassword(passsword, hashedPassword string) (bool, error) {
	parts := strings.Split(hashedPassword, "$")

	if len(parts) != 7 {
		return false, errors.New("invalid hash format")
	}

	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return false, errors.New("invalid hash parameters")
	}

	memVal, err := parseParam(params[0], "m=")
	if err != nil {
		return false, err
	}

	iterVal, err := parseParam(params[1], "t=")
	if err != nil {
		return false, err
	}

	threadsVal, err := parseParam(params[2], "p=")
	if err != nil {
		return false, err
	}

	keyLen, err := parseParam(parts[4], "kl=")
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, errors.New("invalid salt encoding")
	}

	computedHash := argon2.IDKey(
		[]byte(passsword),
		salt,
		uint32(iterVal),
		uint32(memVal),
		uint8(threadsVal),
		uint32(keyLen),
	)

	computedHashEncoded := base64.RawStdEncoding.EncodeToString(computedHash)

	return computedHashEncoded == parts[6], nil
}

func parseParam(param, prefix string) (int, error) {
	if !strings.HasPrefix(param, prefix) {
		return 0, fmt.Errorf("expected parameter prefix %s", prefix)
	}

	val, err := strconv.Atoi(strings.TrimPrefix(param, prefix))
	if err != nil {
		return 0, fmt.Errorf("invalid parameter value: %w", err)
	}

	return val, nil
}
