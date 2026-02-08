package domain

import "errors"

// v0.2: optimistic locking / audit-core
var ErrConflict = errors.New("conflict: base version does not match current head")

// strict audit requirements
var ErrAuditNotConfigured = errors.New("audit log not configured")
var ErrClockNotConfigured = errors.New("clock not configured")
