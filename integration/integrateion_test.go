package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RWCWhen", func() {
	Context("with a country flag", func() {
		It("returns the country", func() {
			By("executing it")
			inputCountry := "Ireland"
			outputCountry, err := runRWCWhen([]string{"--country", inputCountry}, 0)

			Expect(err).ToNot(HaveOccurred())
			Expect(outputCountry).To(ContainSubstring(inputCountry))
		})

	})

	Context("with a group flag", func() {
		It("returns the group", func() {
			By("executing it")
			inputGroup := "A"
			outputGroup, err := runRWCWhen([]string{"--group", inputGroup}, 0)

			Expect(err).ToNot(HaveOccurred())
			Expect(outputGroup).To(ContainSubstring(inputGroup))
		})

	})
})
