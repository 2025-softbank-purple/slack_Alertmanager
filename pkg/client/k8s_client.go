package client

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sClient wraps Kubernetes client operations
type K8sClient struct {
	clientset kubernetes.Interface
}

// NewK8sClient creates a new Kubernetes client
func NewK8sClient() (*K8sClient, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &K8sClient{clientset: clientset}, nil
}

// getKubeConfig retrieves Kubernetes configuration
func getKubeConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to kubeconfig file
	config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	return config, nil
}

// ListNodes lists all nodes in the cluster
func (c *K8sClient) ListNodes() (*corev1.NodeList, error) {
	return c.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

// WatchNodes watches for node changes
func (c *K8sClient) WatchNodes() (watch.Interface, error) {
	return c.clientset.CoreV1().Nodes().Watch(context.TODO(), metav1.ListOptions{})
}

// GetNode retrieves a specific node by name
func (c *K8sClient) GetNode(name string) (*corev1.Node, error) {
	return c.clientset.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
}

