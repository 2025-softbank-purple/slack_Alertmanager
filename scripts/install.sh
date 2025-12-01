#!/bin/bash

set -e

NAMESPACE="${MONITORING_NAMESPACE:-monitoring}"

echo "Installing Kubernetes Auto-Monitoring System..."
echo "Namespace: $NAMESPACE"

# Create namespace if it doesn't exist
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# Install Prometheus Operator
echo "Installing Prometheus Operator..."
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace "$NAMESPACE" \
  --create-namespace \
  --values charts/prometheus-stack/values.yaml \
  --wait --timeout 5m

# Install Grafana
echo "Installing Grafana..."
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install grafana grafana/grafana \
  --namespace "$NAMESPACE" \
  --values charts/grafana/values.yaml \
  --wait --timeout 5m

# Deploy node-exporter DaemonSet (automatically deploys to all nodes)
echo "Deploying node-exporter DaemonSet..."
kubectl apply -f configs/node-exporter/daemonset.yaml

# Create node-exporter Service
echo "Creating node-exporter Service..."
kubectl apply -f configs/node-exporter/service.yaml

# Create ServiceMonitor (automatically discovers all node-exporter services)
echo "Creating ServiceMonitor for node-exporter..."
kubectl apply -f configs/prometheus/servicemonitor.yaml

# Apply PrometheusRule
echo "Applying PrometheusRule..."
kubectl apply -f configs/prometheus/alert-rules.yaml

# Apply Grafana dashboards
echo "Applying Grafana dashboards..."
kubectl apply -f configs/grafana/dashboards/configmap.yaml

echo ""
echo "Installation completed!"
echo ""
echo "To access Grafana, run: kubectl port-forward -n $NAMESPACE svc/grafana 3000:80"
echo "Default credentials: admin / prom-operator"
echo ""
echo "To access Prometheus, run: kubectl port-forward -n $NAMESPACE svc/prometheus-kube-prometheus-prometheus 9090:9090"
echo ""
echo "Note: node-exporter DaemonSet will automatically deploy to all nodes (existing and new)."
echo "      Prometheus will automatically discover all node-exporter services via ServiceMonitor."

