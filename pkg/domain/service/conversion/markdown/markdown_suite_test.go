package markdown

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMarkdown(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Markdown Suite")
}
