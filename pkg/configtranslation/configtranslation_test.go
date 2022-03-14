package configtranslation

import (
	"testing"

	"github.com/SPANDigital/presidium-hugo/pkg/config"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
				jekyllConfig := JekyllConfig{
					Logo: "${baseurl}/media/images/logo.svg",
				}
				logoPrefix := "/images/"
				hugoConfig := ConvertConfig(&jekyllConfig, logoPrefix, map[string]interface{}{})
				Expect(hugoConfig.Params["logo"]).To(Equal("/images/logo.svg"))
			})
		})
		Context("When the configs are not correctly populated", func() {
			It("Should translate to nothing", func() {
				jekyllConfig := JekyllConfig{
					Logo: "nothing/to/see/here",
				}
				logoPrefix := "/images/"
				hugoConfig := ConvertConfig(&jekyllConfig, logoPrefix, map[string]interface{}{})
				Expect(hugoConfig.Params["logo"]).To(Equal(""))
			})
		})
		Context("When brand url is set", func() {
			It("should contain the brand module", func() {
				jekyllConfig := JekyllConfig{}
				config.Flags.BrandTheme = "brand"
				hugoConfig := ConvertConfig(&jekyllConfig, "", map[string]interface{}{})
				Expect(hugoConfig.Module.Imports).To(ContainElement(getImportModule("brand")))
			})
		})
		Context("When markup style is set", func() {
			It("should contain the markup style", func() {
				jekyllConfig := JekyllConfig{}
				config.Flags.Style = "test"
				hugoConfig := ConvertConfig(&jekyllConfig, "", map[string]interface{}{})
				Expect(hugoConfig.Markup.Highlight.Style).To(Equal("test"))
			})
		})
	})
})
