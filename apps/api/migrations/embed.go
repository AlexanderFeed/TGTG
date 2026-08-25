package migrations

// Package migrations contains versioned SQL files embedded in the API binary.
// database.Migrate gives this virtual filesystem to Goose at startup.

import "embed"

// Files contains the versioned SQL migrations compiled into the API binary.
// `go:embed` is a compiler directive, so keep it directly above the variable.
//
//go:embed *.sql
var Files embed.FS
