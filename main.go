package main

import (
	"github.com/SPANDigital/presidium-hugo/cmd"
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/template"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/themes"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
)

func main() {
	template.SetFS(templatesFS)
	if err := model.LoadTemplates(templatesFS); err != nil {
		log.Fatal(err)
	}
	themes.SetZip(themesZip)
	cmd.Execute()
}
