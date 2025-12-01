#!/bin/bash

set -e

echo "=========================================="
echo "DaemonSet 자동 확장 확인"
echo "=========================================="
echo ""

# 현재 노드 수 확인
NODE_COUNT=$(kubectl get nodes --no-headers | wc -l)
echo "1. 현재 클러스터 노드 수: $NODE_COUNT"
kubectl get nodes
echo ""

# DaemonSet 상태 확인
echo "2. DaemonSet 상태:"
kubectl get daemonset node-exporter -n monitoring
echo ""

DESIRED=$(kubectl get daemonset node-exporter -n monitoring -o jsonpath='{.status.desiredNumberScheduled}')
READY=$(kubectl get daemonset node-exporter -n monitoring -o jsonpath='{.status.numberReady}')
AVAILABLE=$(kubectl get daemonset node-exporter -n monitoring -o jsonpath='{.status.numberAvailable}')

echo "   - Desired (원하는 Pod 수): $DESIRED"
echo "   - Ready (준비된 Pod 수): $READY"
echo "   - Available (사용 가능한 Pod 수): $AVAILABLE"
echo ""

# Pod 상태 확인
echo "3. node-exporter Pod 상태 (노드별):"
kubectl get pods -n monitoring -l app=node-exporter -o wide
echo ""

POD_COUNT=$(kubectl get pods -n monitoring -l app=node-exporter --no-headers | wc -l)
echo "   총 Pod 수: $POD_COUNT (노드 수와 일치해야 함: $NODE_COUNT)"
echo ""

# Service Endpoints 확인
echo "4. Service Endpoints (모든 node-exporter 인스턴스):"
kubectl get endpoints -n monitoring node-exporter
echo ""

ENDPOINT_COUNT=$(kubectl get endpoints -n monitoring node-exporter -o jsonpath='{.subsets[0].addresses[*].ip}' | wc -w)
echo "   감지된 Endpoint 수: $ENDPOINT_COUNT"
echo ""

# Prometheus 타겟 확인
echo "5. Prometheus 타겟 확인:"
TARGET_COUNT=$(kubectl exec -n monitoring prometheus-prometheus-kube-prometheus-prometheus-0 -c prometheus -- \
  wget -qO- http://localhost:9090/api/v1/targets 2>/dev/null | \
  grep -o '"job":"node-exporter"' | wc -l 2>/dev/null || echo "0")

if [ "$TARGET_COUNT" -gt "0" ]; then
    echo "   ✓ Prometheus가 $TARGET_COUNT 개의 node-exporter 타겟을 감지했습니다"
else
    echo "   ⚠ Prometheus 타겟 확인 실패 (Prometheus가 아직 시작 중일 수 있음)"
fi
echo ""

# 자동 확장 동작 설명
echo "=========================================="
echo "자동 확장 동작 원리"
echo "=========================================="
echo ""
echo "DaemonSet은 Kubernetes의 기본 기능으로:"
echo "  - 모든 노드에 자동으로 Pod를 배포합니다"
echo "  - 새 노드가 추가되면 자동으로 해당 노드에 Pod를 생성합니다"
echo "  - 노드가 제거되면 해당 노드의 Pod도 자동으로 제거됩니다"
echo ""
echo "현재 상태:"
if [ "$POD_COUNT" -eq "$NODE_COUNT" ] && [ "$DESIRED" -eq "$NODE_COUNT" ]; then
    echo "  ✓ 정상: 모든 노드에 Pod가 배포되어 있습니다"
else
    echo "  ⚠ 주의: Pod 수($POD_COUNT)와 노드 수($NODE_COUNT)가 일치하지 않습니다"
fi
echo ""

# 새 노드 추가 테스트 안내
echo "=========================================="
echo "새 노드 추가 테스트 방법"
echo "=========================================="
echo ""
echo "새 노드를 추가하려면:"
echo ""
echo "방법 1: kind config 파일 사용 (권장)"
echo "  1. kind-multi-node.yaml 파일을 사용하여 클러스터 재생성"
echo "  2. 또는 기존 클러스터에 노드 추가 (kind 버전에 따라 제한적)"
echo ""
echo "방법 2: 수동 확인"
echo "  - 현재 DaemonSet이 $NODE_COUNT 개 노드를 모두 감지하고 있음"
echo "  - 새 노드 추가 시 자동으로 Pod가 생성됨을 확인하려면:"
echo "    watch -n 2 'kubectl get pods -n monitoring -l app=node-exporter -o wide'"
echo ""

