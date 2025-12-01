package helm_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestAlertmanagerDeployment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alertmanager Deployment Suite")
}

var _ = Describe("Alertmanager Deployment", func() {
	var clientset kubernetes.Interface
	var namespace string

	BeforeEach(func() {
		namespace = "monitoring"
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		Expect(err).NotTo(HaveOccurred())
		clientset, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("When Alertmanager is deployed", func() {
		It("should have Alertmanager pods running", func() {
			// This test will fail until Alertmanager is deployed
			pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
				LabelSelector: "app.kubernetes.io/name=alertmanager",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(pods.Items)).To(BeNumerically(">", 0))

			for _, pod := range pods.Items {
				Expect(pod.Status.Phase).To(Equal(corev1.PodRunning))
			}
		})

		It("should have Alertmanager service available", func() {
			// This test will fail until Alertmanager is deployed
			svc, err := clientset.CoreV1().Services(namespace).Get(ctx, "prometheus-kube-prometheus-alertmanager", metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(svc).NotTo(BeNil())
		})
	})
})

var ctx = func() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return ctx
}()

