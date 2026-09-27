package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2id
const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32

	// This gate bounds hashing memory to hashConcurrency x 64 MiB
	// (~256 MiB) — excess logins queue instead of crashing the process.
	hashConcurrency = 4
)

// hashGate is a counting semaphore around Argon2 (see hashConcurrency).
var hashGate = make(chan struct{}, hashConcurrency)

// boundedIDKey runs argon2.IDKey under the concurrency gate.
func boundedIDKey(password, salt []byte, t, mem uint32, p uint8, keyLen uint32) []byte {
	hashGate <- struct{}{}
	defer func() { <-hashGate }()
	return argon2.IDKey(password, salt, t, mem, p, keyLen)
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := boundedIDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func VerifyPassword(password, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}

	var mem, itr uint32
	var par uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &mem, &itr, &par); err != nil {
		return false
	}

	salt, _ := base64.RawStdEncoding.DecodeString(parts[4])
	decodedHash, _ := base64.RawStdEncoding.DecodeString(parts[5])

	compHash := boundedIDKey([]byte(password), salt, itr, mem, par, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(compHash, decodedHash) == 1
}

func GenerateTempPassword(tempPassword string, length int) (string, error) {
	out := make([]byte, length)
	for i := range out {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(tempPassword))))
		if err != nil {
			return "", err
		}
		out[i] = tempPassword[n.Int64()]
	}
	return string(out), nil
}
