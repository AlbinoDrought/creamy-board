package static

import "embed"

//go:embed css/* js/* favicon.ico
var FS embed.FS
