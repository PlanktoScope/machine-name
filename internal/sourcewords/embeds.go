// Package sourcewords provides embedded lists of words which can be used as sources of words to
// compile into wordlists for constructing names.
package sourcewords

import (
	"embed"
)

//go:embed *
var WordlistsFS embed.FS
