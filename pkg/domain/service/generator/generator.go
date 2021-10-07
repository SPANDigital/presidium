package generator

import (
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
)

type Generator interface {
	Run(target model.InitialSiteTarget) error
}
