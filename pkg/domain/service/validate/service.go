package validate

import model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/validate"

type Validator interface {
	Validate() (model.Report, error)
	IsLocal() bool
}
