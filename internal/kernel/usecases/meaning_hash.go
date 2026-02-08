package usecases

import (
	"crypto/sha256"
	"encoding/hex"

	"digiemu-core/internal/kernel/domain"
)

// ComputeMeaningHash canonicalizes the provided Meaning and returns the
// hex-encoded SHA-256 digest. It follows the same minimal canonicalization
// rules as other snapshot artifacts so that hashes are deterministic.
func ComputeMeaningHash(m domain.Meaning) (string, error) {
	canon, err := canonicalJSON(m)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(canon))
	return hex.EncodeToString(sum[:]), nil
}
