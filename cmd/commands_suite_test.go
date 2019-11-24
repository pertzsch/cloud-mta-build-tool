package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var _ = BeforeSuite(func() {
	logs.Logger = logs.NewLogger()
})

func executeAndProvideOutput(execute func() error) (string, error) {
	old := os.Stdout // keep backup of the real stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	err = execute()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			fmt.Println(err)
		}
		outC <- buf.String()
	}()

	os.Stdout = old // restoring the real stdout
	// back to normal state
	_ = w.Close()
	out := <-outC
	return out, err
}

func createFileInTmpFolder(projectName string, path ...string) {
	file, err := os.Create(getFullPathInTmpFolder(projectName, path...))
	Ω(err).Should(Succeed())
	err = file.Close()
	Ω(err).Should(Succeed())
}

func createDirInTmpFolder(projectName string, path ...string) {
	err := dir.CreateDirIfNotExist(getFullPathInTmpFolder(projectName, path...))
	Ω(err).Should(Succeed())
}

func getFullPathInTmpFolder(projectName string, path ...string) string {
	pathWithResultFolder := []string{"result", "." + projectName + "_mta_build_tmp"}
	pathWithResultFolder = append(pathWithResultFolder, path...)
	return getTestPath(pathWithResultFolder...)
}
