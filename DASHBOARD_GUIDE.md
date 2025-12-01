# Grafana 대시보드 확인 가이드

## Grafana 접근 방법

### 1. 포트 포워딩으로 접근 (로컬 테스트)

```bash
# Grafana 포트 포워딩 시작
kubectl port-forward -n monitoring svc/grafana 3000:80
```

**중요**: 이 명령은 포그라운드에서 실행되므로, 별도 터미널 창에서 실행하거나 백그라운드로 실행하세요.

백그라운드 실행:
```bash
kubectl port-forward -n monitoring svc/grafana 3000:80 &
```

### 2. 브라우저에서 접속

포트 포워딩 후 브라우저에서 접속:
- **URL**: http://localhost:3000
- **Username**: `admin`
- **Password**: `prom-operator`

## 대시보드 확인 방법

### 방법 1: 대시보드 브라우저에서 찾기

1. Grafana에 로그인
2. 좌측 메뉴에서 **"Dashboards"** 클릭
3. **"Browse"** 선택
4. 대시보드 목록에서 다음을 찾을 수 있습니다:
   - **Node Exporter Full** (ConfigMap에서 자동 프로비저닝)
   - **Node Exporter for Prometheus** (Grafana 공식 대시보드, ID: 1860)

### 방법 2: 직접 URL로 접근

대시보드 ID를 알고 있다면:
- http://localhost:3000/d/<dashboard-id>

### 방법 3: 검색 기능 사용

1. Grafana 상단의 검색창 클릭
2. "node" 또는 "exporter" 검색
3. 관련 대시보드 선택

## 대시보드 상태 확인

### ConfigMap 확인

```bash
# 대시보드 ConfigMap 확인
kubectl get configmap -n monitoring | grep grafana

# 대시보드 ConfigMap 상세 확인
kubectl get configmap grafana-dashboards -n monitoring -o yaml
```

### Grafana Pod 로그 확인

```bash
# Grafana Pod 이름 확인
GRAFANA_POD=$(kubectl get pods -n monitoring -l app.kubernetes.io/name=grafana -o jsonpath='{.items[0].metadata.name}')

# 로그 확인 (대시보드 로딩 관련)
kubectl logs -n monitoring $GRAFANA_POD | grep -i dashboard
```

### Grafana API로 대시보드 목록 확인

```bash
# Grafana API를 통한 대시보드 목록 확인
kubectl port-forward -n monitoring svc/grafana 3000:80 &
sleep 2

# API로 대시보드 검색 (인증 필요)
curl -u admin:prom-operator http://localhost:3000/api/search?query=node
```

## 설치된 대시보드

### 1. Node Exporter Full (ConfigMap)

- **위치**: `configs/grafana/dashboards/node-exporter.json`
- **자동 프로비저닝**: ConfigMap으로 자동 설치됨
- **라벨**: `grafana_dashboard: "1"`로 자동 감지

### 2. Node Exporter for Prometheus (Grafana 공식)

- **Grafana.com ID**: 1860
- **설치 방법**: Helm values.yaml에서 자동 설치
- **데이터소스**: Prometheus

## 대시보드가 보이지 않는 경우

### 문제 해결 1: ConfigMap 확인

```bash
# ConfigMap이 제대로 생성되었는지 확인
kubectl get configmap grafana-dashboards -n monitoring

# ConfigMap 내용 확인
kubectl describe configmap grafana-dashboards -n monitoring
```

### 문제 해결 2: Grafana 대시보드 프로비저닝 확인

```bash
# Grafana Pod에서 대시보드 디렉토리 확인
GRAFANA_POD=$(kubectl get pods -n monitoring -l app.kubernetes.io/name=grafana -o jsonpath='{.items[0].metadata.name}')
kubectl exec -n monitoring $GRAFANA_POD -- ls -la /var/lib/grafana/dashboards/
```

### 문제 해결 3: 수동으로 대시보드 가져오기

1. Grafana UI 접속
2. 좌측 메뉴에서 **"+"** > **"Import"** 선택
3. **Grafana.com Dashboard** 탭에서 ID 입력: `1860`
4. **Load** 클릭
5. Prometheus 데이터소스 선택
6. **Import** 클릭

## 빠른 확인 스크립트

```bash
#!/bin/bash
echo "=== Grafana 상태 확인 ==="
kubectl get pods -n monitoring -l app.kubernetes.io/name=grafana
echo ""
echo "=== Grafana Service 확인 ==="
kubectl get svc -n monitoring grafana
echo ""
echo "=== 대시보드 ConfigMap 확인 ==="
kubectl get configmap -n monitoring | grep grafana
echo ""
echo "=== 포트 포워딩 시작 ==="
echo "브라우저에서 http://localhost:3000 접속"
echo "로그인: admin / prom-operator"
kubectl port-forward -n monitoring svc/grafana 3000:80
```

## 대시보드에서 확인할 수 있는 메트릭

node-exporter 대시보드에서 다음 메트릭을 확인할 수 있습니다:

- **CPU 사용률**: `100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`
- **메모리 사용률**: `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`
- **디스크 사용률**: `100 - ((node_filesystem_avail_bytes{mountpoint="/"} * 100) / node_filesystem_size_bytes{mountpoint="/"})`
- **네트워크 트래픽**: `rate(node_network_receive_bytes_total[5m])`
- **시스템 부하**: `node_load1`, `node_load5`, `node_load15`

## 참고

- Grafana 공식 문서: https://grafana.com/docs/grafana/latest/dashboards/
- Node Exporter 대시보드: https://grafana.com/grafana/dashboards/1860

