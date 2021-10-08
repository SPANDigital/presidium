package generator

import (
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
)

// SiteGenerator generates presidium site based on a specific initial site model
type SiteGenerator interface {
	Run(target model.InitialSiteTarget) error
}
