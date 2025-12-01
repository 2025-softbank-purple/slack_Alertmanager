package client_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/promethus-example/pkg/client"
)

func TestLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logger Suite")
}

var _ = Describe("Logger", func() {
	Context("When initializing logger", func() {
		It("should create a logger instance", func() {
			logger := client.NewLogger()
			Expect(logger).NotTo(BeNil())
		})

		It("should log info messages", func() {
			logger := client.NewLogger()
			Expect(logger).NotTo(BeNil())
			// This test will fail until logger is implemented
			Expect(logger.Info("test message")).To(Succeed())
		})

		It("should log error messages", func() {
			logger := client.NewLogger()
			Expect(logger).NotTo(BeNil())
			// This test will fail until logger is implemented
			Expect(logger.Error("test error")).To(Succeed())
		})
	})
})

