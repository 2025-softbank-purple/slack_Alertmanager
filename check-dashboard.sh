#!/bin/bash

set -e

echo "=========================================="
echo "Grafana 대시보드 확인 스크립트"
echo "=========================================="
echo ""

# Grafana Pod 확인
echo "1. Grafana Pod 상태:"
GRAFANA_PODS=$(kubectl get pods -n monitoring -l app.kubernetes.io/name=grafana --no-headers 2>/dev/null | wc -l)
if [ "$GRAFANA_PODS" -eq "0" ]; then
    echo "   ⚠ Grafana가 설치되지 않았습니다."
    echo "   설치하려면: ./scripts/install.sh"
    exit 1
else
    kubectl get pods -n monitoring -l app.kubernetes.io/name=grafana
    echo "   ✓ Grafana Pod가 실행 중입니다"
fi
echo ""

# Grafana Service 확인
echo "2. Grafana Service 확인:"
if kubectl get svc -n monitoring grafana &>/dev/null; then
    kubectl get svc -n monitoring grafana
    echo "   ✓ Grafana Service가 있습니다"
else
    echo "   ⚠ Grafana Service를 찾을 수 없습니다"
fi
echo ""

# 대시보드 ConfigMap 확인
echo "3. 대시보드 ConfigMap 확인:"
DASHBOARD_CM=$(kubectl get configmap -n monitoring grafana-dashboards 2>/dev/null || echo "")
if [ -n "$DASHBOARD_CM" ]; then
    kubectl get configmap -n monitoring grafana-dashboards
    echo "   ✓ 대시보드 ConfigMap이 있습니다"
else
    echo "   ⚠ 대시보드 ConfigMap을 찾을 수 없습니다"
    echo "   생성하려면: kubectl apply -f configs/grafana/dashboards/configmap.yaml"
fi
echo ""

# Grafana Pod 내부 대시보드 확인
echo "4. Grafana Pod 내부 대시보드 디렉토리 확인:"
GRAFANA_POD=$(kubectl get pods -n monitoring -l app.kubernetes.io/name=grafana -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
if [ -n "$GRAFANA_POD" ]; then
    echo "   Pod: $GRAFANA_POD"
    kubectl exec -n monitoring $GRAFANA_POD -- ls -la /var/lib/grafana/dashboards/ 2>/dev/null || echo "   대시보드 디렉토리를 확인할 수 없습니다"
else
    echo "   Grafana Pod를 찾을 수 없습니다"
fi
echo ""

# 접근 방법 안내
echo "=========================================="
echo "Grafana 접근 방법"
echo "=========================================="
echo ""
echo "1. 포트 포워딩 시작:"
echo "   kubectl port-forward -n monitoring svc/grafana 3000:80"
echo ""
echo "2. 브라우저에서 접속:"
echo "   URL: http://localhost:3000"
echo "   Username: admin"
echo "   Password: prom-operator"
echo ""
echo "3. 대시보드 찾기:"
echo "   - 좌측 메뉴 > Dashboards > Browse"
echo "   - 'node' 또는 'exporter' 검색"
echo "   - 'Node Exporter Full' 또는 'Node Exporter for Prometheus' 선택"
echo ""

# 포트 포워딩 시작 옵션
read -p "지금 포트 포워딩을 시작하시겠습니까? (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "포트 포워딩 시작 중... (Ctrl+C로 종료)"
    echo "브라우저에서 http://localhost:3000 접속하세요"
    echo ""
    kubectl port-forward -n monitoring svc/grafana 3000:80
fi

