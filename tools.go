//go:build tools
// +build tools

package tools

// Import packages here to ensure `go mod tidy` records their checksums in go.sum.
import (
	_ "github.com/google/uuid"
)
