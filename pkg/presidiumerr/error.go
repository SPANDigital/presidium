package presidiumerr

const (
	InvalidProjectName  = "invalid project name"
	UnsupportedTemplate = "template not supported"
	UnsupportedTheme    = "theme not supported"
	InvalidTitle        = "invalid title"
)

type GenericError struct {
	Code string
}

func (g GenericError) Error() string {
	return g.Code
}
