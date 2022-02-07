package view

import (
	"embed"
	_ "embed"
)

//go:embed index.html global
var EmbedFs embed.FS
