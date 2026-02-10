package usecases

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"digiemu-core/internal/kernel/domain"
)

// CanonicalizeUncertaintyJSON returns canonical JSON bytes for Uncertainty.
func CanonicalizeUncertaintyJSON(b []byte) ([]byte, error) {
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

func ComputeUncertaintyHash(b []byte) (string, error) {
	canon, err := CanonicalizeUncertaintyJSON(b)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canon)
	return hex.EncodeToString(sum[:]), nil
}

func ComputeUncertaintyHashFromStruct(u domain.Uncertainty) (string, error) {
	s, err := canonicalJSON(u)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:]), nil
}
