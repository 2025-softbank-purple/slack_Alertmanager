package client

import "fmt"

// ErrNodeNotFound is returned when a node is not found
type ErrNodeNotFound struct {
	NodeName string
}

func (e *ErrNodeNotFound) Error() string {
	return fmt.Sprintf("node %s not found", e.NodeName)
}

// ErrHelmInstallFailed is returned when Helm installation fails
type ErrHelmInstallFailed struct {
	ChartName string
	Reason    string
}

func (e *ErrHelmInstallFailed) Error() string {
	return fmt.Sprintf("failed to install chart %s: %s", e.ChartName, e.Reason)
}

// ErrResourceCreationFailed is returned when Kubernetes resource creation fails
type ErrResourceCreationFailed struct {
	ResourceType string
	ResourceName string
	Reason       string
}

func (e *ErrResourceCreationFailed) Error() string {
	return fmt.Sprintf("failed to create %s %s: %s", e.ResourceType, e.ResourceName, e.Reason)
}

// WrapError wraps an error with additional context
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

