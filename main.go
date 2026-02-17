package main

import (
	"github.com/SPANDigital/presidium-hugo/cmd"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/template"
)

func main() {
	template.SetFS(templatesFS)
	cmd.Execute()
}
