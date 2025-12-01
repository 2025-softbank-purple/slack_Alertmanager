#!/bin/bash

set -e

NAMESPACE="${MONITORING_NAMESPACE:-monitoring}"

echo "Uninstalling Kubernetes Auto-Monitoring System..."
echo "Namespace: $NAMESPACE"

# Delete node-exporter resources
echo "Deleting node-exporter resources..."
kubectl delete -f configs/node-exporter/ --ignore-not-found=true

# Delete ServiceMonitor
echo "Deleting ServiceMonitor..."
kubectl delete -f configs/prometheus/servicemonitor.yaml --ignore-not-found=true

# Delete PrometheusRule
echo "Deleting PrometheusRule..."
kubectl delete -f configs/prometheus/alert-rules.yaml --ignore-not-found=true

# Uninstall Helm charts
echo "Uninstalling Helm charts..."
helm uninstall prometheus -n "$NAMESPACE" --ignore-not-found=true
helm uninstall grafana -n "$NAMESPACE" --ignore-not-found=true

# Optionally delete namespace (uncomment if desired)
# echo "Deleting namespace..."
# kubectl delete namespace "$NAMESPACE" --ignore-not-found=true

echo "Uninstallation completed!"

