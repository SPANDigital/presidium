package main

import (
	"github.com/SPANDigital/presidium-hugo/cmd"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/template"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/themes"
)

func main() {
	template.SetFS(templatesFS)
	themes.SetFS(themesFS)
	cmd.Execute()
}
