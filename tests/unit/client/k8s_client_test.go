package client_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/promethus-example/pkg/client"
)

func TestK8sClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8sClient Suite")
}

var _ = Describe("K8sClient", func() {
	Context("When creating a Kubernetes client", func() {
		It("should create a client instance", func() {
			client, err := client.NewK8sClient()
			Expect(err).NotTo(HaveOccurred())
			Expect(client).NotTo(BeNil())
		})

		It("should list nodes", func() {
			client, err := client.NewK8sClient()
			Expect(err).NotTo(HaveOccurred())
			// This test will fail until k8s_client is implemented
			nodes, err := client.ListNodes()
			Expect(err).NotTo(HaveOccurred())
			Expect(nodes).NotTo(BeNil())
		})

		It("should watch for node changes", func() {
			client, err := client.NewK8sClient()
			Expect(err).NotTo(HaveOccurred())
			// This test will fail until k8s_client is implemented
			watcher, err := client.WatchNodes()
			Expect(err).NotTo(HaveOccurred())
			Expect(watcher).NotTo(BeNil())
		})
	})
})

