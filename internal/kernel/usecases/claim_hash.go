package usecases

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"digiemu-core/internal/kernel/domain"
)

// CanonicalizeClaimSetJSON takes raw JSON bytes of a ClaimSet and returns the
// minimal canonicalized JSON bytes following the project's canonical rules
// (sorted object keys, preserved array order, minified).
func CanonicalizeClaimSetJSON(b []byte) ([]byte, error) {
	var x any
	if err := json.Unmarshal(b, &x); err != nil {
		return nil, err
	}
	s, err := canonicalJSON(x)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// ComputeClaimSetHash computes SHA-256 hex digest over the canonicalized
// ClaimSet JSON bytes.
func ComputeClaimSetHash(b []byte) (string, error) {
	canon, err := CanonicalizeClaimSetJSON(b)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canon)
	return hex.EncodeToString(sum[:]), nil
}

// convenience wrapper for typed ClaimSet values
func ComputeClaimSetHashFromStruct(cs domain.ClaimSet) (string, error) {
	s, err := canonicalJSON(cs)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:]), nil
}
