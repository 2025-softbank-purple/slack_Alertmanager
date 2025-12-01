#!/bin/bash

set -e

echo "=========================================="
echo "자동 확장 테스트 시작"
echo "=========================================="
echo ""

# 현재 상태 확인
echo "1. 현재 노드 상태:"
kubectl get nodes
echo ""

echo "2. 현재 node-exporter Pod 상태:"
kubectl get pods -n monitoring -l app=node-exporter -o wide
echo ""

# 새 노드 추가
echo "3. 새 노드를 kind 클러스터에 추가 중..."
kind create node --name test-worker --cluster test-monitoring
echo ""

# 노드가 Ready 상태가 될 때까지 대기
echo "4. 새 노드가 Ready 상태가 될 때까지 대기 중..."
sleep 10
kubectl wait --for=condition=Ready node/test-monitoring-test-worker --timeout=60s || true
echo ""

# 노드 상태 확인
echo "5. 업데이트된 노드 상태:"
kubectl get nodes
echo ""

# DaemonSet이 새 노드에 Pod를 배포할 때까지 대기
echo "6. DaemonSet이 새 노드에 Pod를 배포할 때까지 대기 중..."
echo "   (최대 60초 대기)"
for i in {1..12}; do
    NEW_POD_COUNT=$(kubectl get pods -n monitoring -l app=node-exporter -o wide | grep -c test-worker || echo "0")
    if [ "$NEW_POD_COUNT" -ge "1" ]; then
        echo "   ✓ 새 노드에 Pod가 배포되었습니다!"
        break
    fi
    echo "   대기 중... ($i/12)"
    sleep 5
done
echo ""

# 최종 Pod 상태 확인
echo "7. 최종 node-exporter Pod 상태:"
kubectl get pods -n monitoring -l app=node-exporter -o wide
echo ""

# Service Endpoints 확인
echo "8. Service Endpoints 확인 (모든 node-exporter 인스턴스):"
kubectl get endpoints -n monitoring node-exporter
echo ""

# Prometheus 타겟 확인 (선택적)
echo "9. Prometheus 타겟 확인:"
echo "   다음 명령으로 Prometheus에서 타겟을 확인할 수 있습니다:"
echo "   kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090"
echo "   브라우저에서 http://localhost:9090/targets 접속"
echo ""

# 정리 안내
echo "=========================================="
echo "테스트 완료!"
echo "=========================================="
echo ""
echo "정리하려면:"
echo "  kind delete node test-monitoring-test-worker --cluster test-monitoring"

