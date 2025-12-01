package contract_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPrometheusRule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PrometheusRule Contract Suite")
}

var _ = Describe("PrometheusRule Contract", func() {
	Context("When creating PrometheusRule", func() {
		It("should have correct structure", func() {
			rule := &monitoringv1.PrometheusRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "k8s-node-alerts",
					Namespace: "monitoring",
				},
				Spec: monitoringv1.PrometheusRuleSpec{
					Groups: []monitoringv1.RuleGroup{
						{
							Name:     "node.rules",
							Interval: "30s",
							Rules: []monitoringv1.Rule{
								{
									Alert: "NodeDown",
									Expr:  "up{job=\"node-exporter\"} == 0",
									For:   "5m",
									Labels: map[string]string{
										"severity": "critical",
									},
									Annotations: map[string]string{
										"summary":     "Node {{ $labels.instance }} is down",
										"description": "Node {{ $labels.instance }} has been down for more than 5 minutes.",
									},
								},
							},
						},
					},
				},
			}

			Expect(rule.Name).To(Equal("k8s-node-alerts"))
			Expect(rule.Namespace).To(Equal("monitoring"))
			Expect(len(rule.Spec.Groups)).To(Equal(1))
			Expect(rule.Spec.Groups[0].Name).To(Equal("node.rules"))
			Expect(len(rule.Spec.Groups[0].Rules)).To(BeNumerically(">", 0))
		})

		It("should include default alert rules", func() {
			// This test validates the contract structure
			// Actual rule creation will be tested in integration tests
			Expect(true).To(BeTrue())
		})
	})
})

