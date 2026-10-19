// Package exercises embeds the curriculum so release binaries work from any directory.
package exercises

import "embed"

// Extensions are listed explicitly so stray __pycache__ directories are never embedded.
//
//go:embed modules.toml */*/*.toml */*/*.py */*/*.md
var FS embed.FS
