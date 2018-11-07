package integration_test

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	var (
		gsaCmd    *exec.Cmd
		grootfs   string
		config    string
		gsaOutput *bytes.Buffer
	)

	JustBeforeEach(func() {
		gsaCmd = exec.Command(gsaBin, "--grootfs-bin", grootfs, "--grootfs-config", config)
		gsaOutput = bytes.NewBuffer([]byte{})
		gsaCmd.Stdout = gsaOutput
		gsaCmd.Stderr = gsaOutput
	})

	BeforeEach(func() {
		grootfs = fakeGrootBin
		config = grootConfigFile
	})

	It("gets total disk usage of all containers in depot", func() {
		err := gsaCmd.Run()
		Expect(err).NotTo(HaveOccurred())
		Expect(gsaOutput.String()).To(ContainSubstring("\"total_bytes_containers\":246"))
	})

	It("gets total disk usage of all layers in store", func() {
		err := gsaCmd.Run()
		Expect(err).NotTo(HaveOccurred())
		Expect(gsaOutput.String()).To(ContainSubstring("\"total_bytes_layers\":2468"))
	})

	It("gets total disk usage of all active layers", func() {
		err := gsaCmd.Run()
		Expect(err).NotTo(HaveOccurred())
		Expect(gsaOutput.String()).To(ContainSubstring("\"total_bytes_active_layers\":1234"))
	})

	It("returns total store usage (excl meta)", func() {
		err := gsaCmd.Run()
		Expect(err).NotTo(HaveOccurred())
		Expect(gsaOutput.String()).To(ContainSubstring("\"total_bytes_store\":2714"))
	})

	Context("the -h flag", func() {
		JustBeforeEach(func() {
			gsaCmd.Args = append(gsaCmd.Args, "-r")
		})

		It("returns a human readable result", func() {
			err := gsaCmd.Run()
			Expect(err).NotTo(HaveOccurred())
			Expect(gsaOutput.String()).To(Equal("Containers: 246 B\nLayers: 2.5 kB (of which Active: 1.2 kB)\n2.7 kB\n"))
		})
	})

	Context("when the grootfs binary cannot be found", func() {
		BeforeEach(func() {
			grootfs = "/not/found"
		})

		It("fails", func() {
			err := gsaCmd.Run()
			Expect(err).To(HaveOccurred())
			Expect(gsaOutput.String()).To(ContainSubstring("grootfs not found"))
		})
	})

	Context("when groot's config cannot be found", func() {
		BeforeEach(func() {
			config = "/not/found"
		})

		It("Fails", func() {
			err := gsaCmd.Run()
			Expect(err).To(HaveOccurred())
			Expect(gsaOutput.String()).To(ContainSubstring("grootfs config not found"))
		})
	})

	Context("when the grootfs store cannot be parsed from the config", func() {
		BeforeEach(func() {
			config = filepath.Join(grootConfigPath, "garbage-config.yml")
			Expect(ioutil.WriteFile(config, []byte("some-garbage"), 0755)).To(Succeed())
		})

		It("fails", func() {
			err := gsaCmd.Run()
			Expect(err).To(HaveOccurred())
			Expect(gsaOutput.String()).To(ContainSubstring("unmarshal errors"))
		})
	})
})
