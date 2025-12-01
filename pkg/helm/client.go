package helm

import (
	"context"
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/client-go/rest"
)

// HelmClient wraps Helm action operations
type HelmClient struct {
	settings *cli.EnvSettings
	config   *action.Configuration
}

// NewHelmClient creates a new Helm client
func NewHelmClient(namespace string) (*HelmClient, error) {
	settings := cli.New()
	settings.SetNamespace(namespace)

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secret", func(format string, v ...interface{}) {
		// Log function - can be customized
	}); err != nil {
		return nil, fmt.Errorf("failed to initialize Helm action config: %w", err)
	}

	return &HelmClient{
		settings: settings,
		config:   actionConfig,
	}, nil
}

// GetInstallAction returns a Helm install action
func (c *HelmClient) GetInstallAction() *action.Install {
	installAction := action.NewInstall(c.config)
	installAction.Namespace = c.settings.Namespace()
	installAction.CreateNamespace = true
	return installAction
}

// GetUpgradeAction returns a Helm upgrade action
func (c *HelmClient) GetUpgradeAction() *action.Upgrade {
	upgradeAction := action.NewUpgrade(c.config)
	upgradeAction.Namespace = c.settings.Namespace()
	return upgradeAction
}

// GetUninstallAction returns a Helm uninstall action
func (c *HelmClient) GetUninstallAction() *action.Uninstall {
	uninstallAction := action.NewUninstall(c.config)
	return uninstallAction
}

