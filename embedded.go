package main

import "embed"

//go:embed all:templates
var templatesFS embed.FS

//go:embed all:themes
var themesFS embed.FS
