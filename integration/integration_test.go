package integration_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RWCWhen", func() {
	Context("with a country flag", func() {
		It("returns the country", func() {
			By("executing it")
			inputCountry := "Ireland"
			outputCountry, _ := runRWCWhen([]string{"--country", inputCountry}, 0)
			outpurCountryString := strings.TrimSpace(string(getContentsOfReader(outputCountry)))
			Expect(outpurCountryString).To(ContainSubstring(inputCountry))
		})

	})
})
