package migrations

import "embed"

//go:embed *.sql fs/*
var FS embed.FS
