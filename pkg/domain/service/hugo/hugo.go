package hugo

import "github.com/gohugoio/hugo/commands"

type Service struct {
}

func New() Service {
	return Service{}
}

func (s Service) Execute(args ...string) {
	commands.Execute(args)
}
