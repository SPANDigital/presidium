package utils

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = Describe("utils", func() {
	When("deriving article title", func() {
		givenExpectations := map[string]string{
			"0.1.1-sales-exercises": "Sales Exercises",
			"financial--activities": "Financial Activities",
			"1.1.1-the-happy-path":  "The Happy Path",
		}
		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("title of \"%s\" should be: \"%s\"", given, expecting)
			a, b := given, expecting
			It(should, func() {
				actual := UnSlugify(a)
				Expect(actual).Should(Equal(b))
			})
		}
	})

	When("deriving article slug", func() {
		givenExpectations := map[string]string{
			"v0 .18.6_8.":      "v0-18-6-8",
			"update_terraform": "update-terraform",
		}
		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("slug of \"%s\" should be: \"%s\"", given, expecting)
			a, b := given, expecting
			It(should, func() {
				actual := Slugify(a)
				Expect(actual).Should(Equal(b))
			})
		}
	})

	When("deriving slug from title", func() {
		givenExpectations := map[string]string{
			"Troubleshooting":                    "troubleshooting",
			"Set up Development Environment":     "set-up-development-environment",
			"\"Set up Development Environment\"": "set-up-development-environment",
			"Introduction & Overview":            "introduction-and-overview",
		}

		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("title of \"%s\" should be: \"%s\"", given, expecting)
			a, b := given, expecting
			It(should, func() {
				actual := TitleToSlug(a)
				Expect(actual).Should(Equal(b))
			})
		}
	})
})
