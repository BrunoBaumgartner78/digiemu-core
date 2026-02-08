package domain

import "errors"

// v0.2.2.1: head pointer exists but cannot be resolved in versions list
var ErrInconsistentHead = errors.New("inconsistent head version pointer")
