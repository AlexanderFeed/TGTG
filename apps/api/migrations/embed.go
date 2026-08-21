package migrations

import "embed"

// Files contains the versioned SQL migrations compiled into the API binary.
//
//go:embed *.sql
var Files embed.FS
