package configtranslation_test

import (
	"github.com/SPANDigital/presidium-hugo/pkg/configtranslation"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestConfigTranslation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configuration Translation Suite")
}

var _ = Describe("ConfigTranslation", func() {
	BeforeSuite(func() {
		colors.Setup()
	})
	Describe("Converting from Jekyll config to Hugo Config", func() {
		Context("When the configs are correctly populated", func() {
			It("Should translate the path for the Logo correctly", func() {
				jekyllConfig := configtranslation.JekyllConfig{
					Logo: "${baseurl}/media/images/logo.svg",
				}
				logoPrefix := "/images/"
				hugoConfig := configtranslation.ConvertConfig(&jekyllConfig, logoPrefix, map[string]interface{}{})
				Expect(hugoConfig.Params["logo"]).To(Equal("/images/logo.svg"))
			})
		})
		Context("When the configs are not correctly populated", func() {
			It("Should translate to nothing", func() {
				jekyllConfig := configtranslation.JekyllConfig{
					Logo: "nothing/to/see/here",
				}
				logoPrefix := "/images/"
				hugoConfig := configtranslation.ConvertConfig(&jekyllConfig, logoPrefix, map[string]interface{}{})
				Expect(hugoConfig.Params["logo"]).To(Equal(""))
			})
		})
	})

})
