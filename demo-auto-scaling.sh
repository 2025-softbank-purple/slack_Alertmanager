#!/bin/bash

set -e

echo "=========================================="
echo "자동 확장 데모 스크립트"
echo "=========================================="
echo ""
echo "이 스크립트는:"
echo "  1. 현재 클러스터를 삭제"
echo "  2. 3개 노드(1 control-plane + 2 worker)로 클러스터 재생성"
echo "  3. 모니터링 시스템 재설치"
echo "  4. DaemonSet이 모든 노드에 자동 배포되는지 확인"
echo ""
read -p "계속하시겠습니까? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "취소되었습니다."
    exit 0
fi

echo ""
echo "=== 1단계: 기존 클러스터 삭제 ==="
kind delete cluster --name test-monitoring 2>/dev/null || echo "클러스터가 없거나 이미 삭제됨"
echo ""

echo "=== 2단계: 여러 노드로 클러스터 생성 ==="
kind create cluster --config kind-multi-node.yaml
echo ""

echo "=== 3단계: 노드 확인 ==="
kubectl get nodes
echo ""

echo "=== 4단계: 모니터링 시스템 설치 ==="
./scripts/install.sh
echo ""

echo "=== 5단계: DaemonSet 자동 배포 확인 (30초 대기) ==="
sleep 30
echo ""

echo "=== 6단계: 결과 확인 ==="
NODE_COUNT=$(kubectl get nodes --no-headers | wc -l)
POD_COUNT=$(kubectl get pods -n monitoring -l app=node-exporter --no-headers | wc -l)

echo "노드 수: $NODE_COUNT"
echo "node-exporter Pod 수: $POD_COUNT"
echo ""

if [ "$POD_COUNT" -eq "$NODE_COUNT" ]; then
    echo "✓ 성공! 모든 노드에 Pod가 자동으로 배포되었습니다!"
else
    echo "⚠ 주의: Pod 수가 노드 수와 일치하지 않습니다. 잠시 후 다시 확인해주세요."
fi
echo ""

echo "=== Pod 상세 정보 ==="
kubectl get pods -n monitoring -l app=node-exporter -o wide
echo ""

echo "=== Service Endpoints ==="
kubectl get endpoints -n monitoring node-exporter
echo ""

echo "=========================================="
echo "데모 완료!"
echo "=========================================="
echo ""
echo "추가 확인:"
echo "  ./test-daemonset-scaling.sh  # 상세 상태 확인"
echo "  kubectl get daemonset -n monitoring  # DaemonSet 상태"

