package impl

import (
	"errors"
	"fmt"
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/template"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"path/filepath"
)

type gen struct {
	filesystem.FileSystem
	template.Service
}

func (g *gen) Run(target model.InitialSiteTarget) error {
	if err := g.prepare(target); err != nil {
		return err
	}
	return g.processTemplates(target)
}

func (g *gen) processTemplates(target model.InitialSiteTarget) error {

	or := func(s1, s2 string) string {
		if len(s1) > 0 {
			return s1
		} else {
			return s2
		}
	}

	return g.ProcessDirTemplates(target.Template.Code(), target.SiteTargetDirectory, model.TemplateParameters{
		Title:       or(target.SiteTitle, target.SiteName),
		ProjectName: target.SiteName,
		Theme:       target.Theme.ModulePath(),
		Template:    target.Template.Code(),
		Brand:       target.BrandingModelUrl,
	})
}

func (g gen) prepare(target model.InitialSiteTarget) error {

	dirExists := g.DirExists(target.SiteTargetDirectory)

	if !dirExists {
		if err := g.MakeDirs(target.SiteTargetDirectory); err != nil {
			return err
		}
	} else {
		switch target.WhenSiteExists {
		case model.AbortWhenTargetSiteExists:
			return errors.New(fmt.Sprintf("site already exists here: %s", target.SiteTargetDirectory))
		case model.ReplaceTargetSiteIfExists:
			if err := g.EmptyDir(target.SiteTargetDirectory); err != nil {
				return err
			}
		}
	}

	if err := g.MakeDirs(filepath.Join("static")); err != nil {
		return err
	}

	return nil
}

func New() generator.Generator {
	return &gen{
		FileSystem: filesystem.New(),
		Service:    template.New(),
	}
}
