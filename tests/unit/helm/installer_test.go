package helm_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/promethus-example/pkg/helm"
)

func TestInstaller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helm Installer Suite")
}

var _ = Describe("Helm Installer", func() {
	Context("When installing Prometheus Operator", func() {
		It("should install Prometheus Operator chart", func() {
			installer := helm.NewInstaller("monitoring")
			Expect(installer).NotTo(BeNil())
			// This test will fail until installer is implemented
			err := installer.InstallPrometheusOperator("prometheus-community/kube-prometheus-stack", "55.0.0")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When installing Grafana", func() {
		It("should install Grafana chart", func() {
			installer := helm.NewInstaller("monitoring")
			Expect(installer).NotTo(BeNil())
			// This test will fail until installer is implemented
			err := installer.InstallGrafana("grafana/grafana", "6.50.0")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When installing Alertmanager", func() {
		It("should install Alertmanager via Prometheus Operator", func() {
			installer := helm.NewInstaller("monitoring")
			Expect(installer).NotTo(BeNil())
			// This test will fail until installer is implemented
			err := installer.InstallAlertmanager()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

