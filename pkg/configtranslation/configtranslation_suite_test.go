package configtranslation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfigtranslation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configtranslation Suite")
}
