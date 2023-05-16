// Package sources provides embedded lists of words which can be used as sources for compiling
// wordlists for constructing names.
package sources

import (
	"embed"
)

//go:embed *
var WordlistsFS embed.FS
