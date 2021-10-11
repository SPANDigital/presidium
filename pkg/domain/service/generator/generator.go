package generator

import (
	"errors"
	"fmt"
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/template"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
)

// SiteGenerator generates presidium site based on a specific initial site model
type SiteGenerator interface {
	Run(target model.InitialSiteTarget) error
}

func New() SiteGenerator {
	return &gen{
		FileSystem: filesystem.New(),
		Service:    template.New(),
	}
}

type gen struct {
	filesystem.FileSystem
	template.Service
}

func (g *gen) Run(target model.InitialSiteTarget) error {
	if err := g.prepareSiteTarget(target); err != nil {
		return err
	}
	return g.processTemplates(target)
}

func (g *gen) processTemplates(target model.InitialSiteTarget) error {
	return g.ProcessDirTemplates(
		target.Template.Code(),
		target.SiteTargetDirectory,
		target.GetTemplateParameters(),
	)
}

func (g gen) prepareSiteTarget(t model.InitialSiteTarget) error {

	dirExists := g.DirExists(t.SiteTargetDirectory)

	if !dirExists {
		if err := g.MakeDirs(t.SiteTargetDirectory); err != nil {
			return err
		}
	} else {
		switch t.WhenSiteExists {
		case model.AbortWhenTargetSiteExists:
			return errors.New(fmt.Sprintf("site already exists here: %s", t.SiteTargetDirectory))
		case model.ReplaceTargetSiteIfExists:
			if err := g.EmptyDir(t.SiteTargetDirectory); err != nil {
				return err
			}
		}
	}

	if err := g.MakeDirs(t.AssetsDir()); err != nil {
		return err
	}

	return nil
}
