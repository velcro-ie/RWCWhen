package integration_test

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RWCWhen", func() {
	Context("with a country flag", func() {
		It("returns the country", func() {
			By("executing it")
			inputCountry := "Ireland"
			outputCountry, err := runRWCWHen([]string{"--country", inputCountry}, 0)
			
			Expect(err).ToNot(HaveOccurred())
			Expect(outputCountry).To(ContainSubstring(inputCountry))
		})

	})

	Context("with a group flag", func() {
		It("returns the group", func() {
			By("executing it")
			inputGroup := "A"
			outputGroup, stdErr := runRWCWHen([]string{"--group", inputGroup})
			
			Expect(err).ToNot(HaveOccurred())
			Expect(outputGroup).To(ContainSubstring(inputGroup))
		})

	})
})


func runRWCWhen(args []string, expErrCode int) (stdOut *bytes.Reader, stdErr *bytes.Reader) {
	stdOutBuffer := bytes.Buffer{}
	stdErrBuffer := bytes.Buffer{}

	cmd := exec.Command(pathToBin, args...)

	session, err := gexec.Start(cmd, &stdOutBuffer, &stdErrBuffer)
	Expect(err).ToNot(HaveOccurred())
	<-session.Exited

	stdOut = bytes.NewReader(stdOutBuffer.Bytes())
	stdErr = bytes.NewReader(stdErrBuffer.Bytes())

	if os.Getenv("DEBUG") != "" {
		io.Copy(os.Stdout, stdOut)
		io.Copy(os.Stdout, stdErr)
		stdOut.Seek(0, 0)
		stdErr.Seek(0, 0)
	}

	Eventually(session, time.Minute).Should(gexec.Exit(expErrCode))

	return stdOut, stdErr
}