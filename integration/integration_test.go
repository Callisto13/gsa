package integration_test

import (
	"bytes"
	"os/exec"

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
})
