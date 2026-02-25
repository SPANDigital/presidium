package generator

import (
	"fmt"

	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/template"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
)

// SiteGenerator generates presidium site based on a specific initial site model
type SiteGenerator interface {
	Run(target model.InitialSiteTarget) error
}

func New() (SiteGenerator, error) {
	tmplSvc, err := template.New()
	if err != nil {
		return nil, fmt.Errorf("initializing template service: %w", err)
	}
	return &gen{
		FsUtil:  filesystem.New(),
		Service: tmplSvc,
	}, nil
}

type gen struct {
	filesystem.FsUtil
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
			return fmt.Errorf("site already exists here: %s", t.SiteTargetDirectory)
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
