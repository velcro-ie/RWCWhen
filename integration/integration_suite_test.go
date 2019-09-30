package integration_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	pathToBin                 string
	commitHash, pathToGitRepo string
)

func TestRWCWhen(t *testing.T) {
	RegisterFailHandler(Fail)

	AfterSuite(func() {
		os.RemoveAll(pathToGitRepo)
		gexec.Kill()
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "deplab Suite")
}

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

func getContentsOfReader(r io.Reader) []byte {
	contents, err := ioutil.ReadAll(r)
	Expect(err).NotTo(HaveOccurred())

	return contents
}
