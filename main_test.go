package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var (
	cliPath string
	args    []string
	session *Session
	err     error
)

func getPaths(dirName string) ([]string, error) {
	paths := make([]string, 1)
	err = filepath.Walk(dirName, func(path string, _ os.FileInfo, _ error) (errVal error) {
		subPath := strings.Replace(path, dirName, "", 1)
		paths = append(paths, subPath)
		return errVal
	})
	return paths, err
}

var _ = Describe("hugo-parser", func() {
	BeforeSuite(func() {
		cliPath, err = Build("github.com/EngineerBetter/hugo-parser")
		Ω(err).ShouldNot(HaveOccurred(), "Error building source")
	})

	AfterSuite(func() {
		CleanupBuildArtifacts()
	})

	JustBeforeEach(func() {
		command := exec.Command(cliPath, args...)
		session, err = Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred(), "Error running CLI: "+cliPath)
	})

	Context("when the files exist", func() {
		var dir string
		var targetContents []string

		BeforeEach(func() {
			dir, err = ioutil.TempDir("", "example")
			defer os.RemoveAll(dir)
			Ω(err).ShouldNot(HaveOccurred(), "Error creating temp dir")
			args = []string{"fixtures/hugo-structure.yml", "fixtures/source", dir}
			targetContents, err = getPaths("fixtures/target")
			Ω(err).ShouldNot(HaveOccurred(), "Error checking contents of target")
		})

		It("it reads the file", func() {
			Eventually(session).Should(Exit(0))
			result := `mappings:
- name: A
  exercises:
  - B
  - C
- name: D
  exercises:
  - E`
			Ω(session.Out).Should(Say(result))
			actualContents, err := getPaths(dir)
			Ω(err).ShouldNot(HaveOccurred(), "Error checking contents of actual")
			fmt.Println(targetContents)
			fmt.Println(actualContents)
			Ω(reflect.DeepEqual(targetContents, actualContents)).Should(BeTrue())
		})
	})
})
