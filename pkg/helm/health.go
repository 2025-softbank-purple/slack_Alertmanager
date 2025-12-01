package helm

import (
	"context"
	"fmt"
	"time"

	"github.com/promethus-example/pkg/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// HealthChecker checks the health of deployed components
type HealthChecker struct {
	clientset kubernetes.Interface
	logger    *client.Logger
	namespace string
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(clientset kubernetes.Interface, namespace string) *HealthChecker {
	return &HealthChecker{
		clientset: clientset,
		logger:    client.NewLogger(),
		namespace: namespace,
	}
}

// CheckComponentHealth checks if a component is healthy
func (h *HealthChecker) CheckComponentHealth(componentName string, selector string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pods, err := h.clientset.CoreV1().Pods(h.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return fmt.Errorf("failed to list pods for %s: %w", componentName, err)
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for %s", componentName)
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			return fmt.Errorf("pod %s is not running (status: %s)", pod.Name, pod.Status.Phase)
		}
	}

	h.logger.Info(fmt.Sprintf("Component %s is healthy", componentName))
	return nil
}

// WaitForComponentReady waits for a component to be ready
func (h *HealthChecker) WaitForComponentReady(componentName string, selector string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := h.CheckComponentHealth(componentName, selector); err == nil {
				return nil
			}
		case <-time.After(time.Until(deadline)):
			return fmt.Errorf("timeout waiting for %s to be ready", componentName)
		}
	}
}

