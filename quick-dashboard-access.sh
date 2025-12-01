#!/bin/bash

echo "=========================================="
echo "Grafana 대시보드 빠른 접근"
echo "=========================================="
echo ""
echo "포트 포워딩을 시작합니다..."
echo "브라우저에서 http://localhost:3000 접속하세요"
echo ""
echo "로그인 정보:"
echo "  Username: admin"
echo "  Password: prom-operator"
echo ""
echo "대시보드 찾기:"
echo "  1. 좌측 메뉴에서 'Dashboards' 클릭"
echo "  2. 'Browse' 선택"
echo "  3. 검색창에 'node' 입력"
echo "  4. 'Node Exporter Full' 또는 'Node Exporter for Prometheus' 선택"
echo ""
echo "종료하려면 Ctrl+C를 누르세요"
echo "=========================================="
echo ""

kubectl port-forward -n monitoring svc/grafana 3000:80

