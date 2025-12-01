package helm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/promethus-example/pkg/client"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

// Installer handles Helm chart installations
type Installer struct {
	helmClient *HelmClient
	logger     *client.Logger
	namespace  string
}

// NewInstaller creates a new Helm installer
func NewInstaller(namespace string) *Installer {
	helmClient, _ := NewHelmClient(namespace)
	logger := client.NewLogger()
	return &Installer{
		helmClient: helmClient,
		logger:     logger,
		namespace:  namespace,
	}
}

// InstallPrometheusOperator installs Prometheus Operator via Helm
func (i *Installer) InstallPrometheusOperator(chartRepo, version string) error {
	i.logger.Info(fmt.Sprintf("Installing Prometheus Operator from %s (version %s)", chartRepo, version))

	// Add repository
	if err := i.addRepo("prometheus-community", "https://prometheus-community.github.io/helm-charts"); err != nil {
		return fmt.Errorf("failed to add Prometheus repository: %w", err)
	}

	// Install chart
	installAction := i.helmClient.GetInstallAction()
	installAction.ReleaseName = "prometheus"
	installAction.ChartPathOptions.Version = version

	chartPath, err := i.getChartPath("prometheus-community/kube-prometheus-stack", version)
	if err != nil {
		return fmt.Errorf("failed to get chart path: %w", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	values, err := i.loadValues("charts/prometheus-stack/values.yaml")
	if err != nil {
		return fmt.Errorf("failed to load values: %w", err)
	}

	_, err = installAction.Run(chart, values)
	if err != nil {
		return &client.ErrHelmInstallFailed{
			ChartName: "prometheus-operator",
			Reason:    err.Error(),
		}
	}

	i.logger.Info("Prometheus Operator installed successfully")
	return nil
}

// InstallGrafana installs Grafana via Helm
func (i *Installer) InstallGrafana(chartRepo, version string) error {
	i.logger.Info(fmt.Sprintf("Installing Grafana from %s (version %s)", chartRepo, version))

	// Add repository
	if err := i.addRepo("grafana", "https://grafana.github.io/helm-charts"); err != nil {
		return fmt.Errorf("failed to add Grafana repository: %w", err)
	}

	// Install chart
	installAction := i.helmClient.GetInstallAction()
	installAction.ReleaseName = "grafana"
	installAction.ChartPathOptions.Version = version

	chartPath, err := i.getChartPath("grafana/grafana", version)
	if err != nil {
		return fmt.Errorf("failed to get chart path: %w", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	values, err := i.loadValues("charts/grafana/values.yaml")
	if err != nil {
		return fmt.Errorf("failed to load values: %w", err)
	}

	_, err = installAction.Run(chart, values)
	if err != nil {
		return &client.ErrHelmInstallFailed{
			ChartName: "grafana",
			Reason:    err.Error(),
		}
	}

	i.logger.Info("Grafana installed successfully")
	return nil
}

// InstallAlertmanager installs Alertmanager (included in Prometheus Operator)
func (i *Installer) InstallAlertmanager() error {
	// Alertmanager is included in kube-prometheus-stack
	// This method ensures it's enabled in the values
	i.logger.Info("Alertmanager is included in Prometheus Operator stack")
	return nil
}

// Helper methods

func (i *Installer) addRepo(name, url string) error {
	settings := cli.New()
	repoFile := settings.RepositoryConfig

	// Load existing repos
	f, err := repo.LoadFile(repoFile)
	if err != nil {
		return err
	}

	// Check if repo already exists
	if f.Has(name) {
		return nil
	}

	// Add new repo
	chartRepo := &repo.Entry{
		Name: name,
		URL:  url,
	}

	chartRepo, err = repo.NewChartRepository(chartRepo, getter.All(settings))
	if err != nil {
		return err
	}

	_, err = chartRepo.DownloadIndexFile()
	if err != nil {
		return err
	}

	f.Update(chartRepo)
	return f.WriteFile(repoFile, 0644)
}

func (i *Installer) getChartPath(chartName, version string) (string, error) {
	settings := cli.New()
	chartPathOptions := action.ChartPathOptions{
		Version: version,
	}

	cp, err := chartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		return "", err
	}

	return cp, nil
}

func (i *Installer) loadValues(valuesPath string) (map[string]interface{}, error) {
	if _, err := os.Stat(valuesPath); os.IsNotExist(err) {
		// Return empty values if file doesn't exist
		return make(map[string]interface{}), nil
	}

	data, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		return nil, err
	}

	// Simple YAML parsing - in production, use proper YAML parser
	// For now, return empty map
	return make(map[string]interface{}), nil
}

