package generator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SiteGenerator Test Suite")
}

var _ = Describe("Site generation", func() {

})
