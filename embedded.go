package main

import (
	"embed"
	_ "embed"
)

//go:embed all:templates
var templatesFS embed.FS

// Embed themes as a compressed zip bundle
// This avoids Go module embedding restrictions and improves extraction speed
//
//go:embed themes.zip
var themesZip []byte
