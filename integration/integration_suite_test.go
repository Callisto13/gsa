package integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	tmpDir          string
	gsaBin          string
	fakeGrootBin    string
	gardenDepotDir  string
	grootConfigFile string
	grootStorePath  string
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var err error
		gsaBin, err = gexec.Build("github.com/Callisto13/gsa/cmd")
		Expect(err).NotTo(HaveOccurred())

		fakeGrootBin, err = gexec.Build("github.com/Callisto13/gsa/integration/cmd/fakegrootfs")
		Expect(err).NotTo(HaveOccurred())

		tmpDir, err = ioutil.TempDir("", "gsa-integration")
		Expect(err).NotTo(HaveOccurred())

		grootConfigPath := filepath.Join(tmpDir, "grootfs")
		Expect(os.MkdirAll(grootConfigPath, 0755)).To(Succeed())
		grootStorePath = filepath.Join(tmpDir, "store")
		Expect(os.MkdirAll(grootStorePath, 0755)).To(Succeed())
		grootConfigFile = filepath.Join(grootConfigPath, "config.yml")

		grootImageDir := filepath.Join(grootStorePath, "images")
		Expect(os.MkdirAll(grootImageDir, 0755)).To(Succeed())
		Expect(os.Mkdir(filepath.Join(grootImageDir, "one-image"), 0755)).To(Succeed())
		Expect(os.Mkdir(filepath.Join(grootImageDir, "two-image"), 0755)).To(Succeed())

		grootMetaPath := filepath.Join(grootStorePath, "meta")
		Expect(os.MkdirAll(grootMetaPath, 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(grootMetaPath, "volume-one-volume"), []byte("{\"Size\": 1234}"), 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(grootMetaPath, "volume-two-volume"), []byte("{\"Size\": 1234}"), 0755)).To(Succeed())

		grootDepsPath := filepath.Join(grootMetaPath, "dependencies")
		Expect(os.MkdirAll(grootDepsPath, 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(grootDepsPath, "image:one-image.json"), []byte("[\"one-volume\"]"), 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(grootDepsPath, "image:two-image.json"), []byte("[\"one-volume\"]"), 0755)).To(Succeed())

		type config struct {
			StorePath string `yaml:"store"`
		}
		y := config{StorePath: grootStorePath}
		yml, err := yaml.Marshal(y)
		Expect(err).NotTo(HaveOccurred())
		Expect(ioutil.WriteFile(grootConfigFile, yml, 0755)).To(Succeed())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	RunSpecs(t, "Integration Suite")
}
