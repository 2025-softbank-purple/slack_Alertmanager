# 설치 검증 가이드

## ✅ 현재 상태 확인

설치가 성공적으로 완료되었습니다! 다음 명령으로 상태를 확인할 수 있습니다:

### 1. Pod 상태 확인
```bash
kubectl get pods -n monitoring
```

**예상 결과**: 모든 Pod가 `Running` 상태여야 합니다.

### 2. DaemonSet 확인
```bash
kubectl get daemonset -n monitoring
```

**예상 결과**: `node-exporter` DaemonSet이 모든 노드에 배포되어야 합니다.

### 3. ServiceMonitor 확인
```bash
kubectl get servicemonitor -n monitoring
```

**예상 결과**: `node-exporter` ServiceMonitor가 생성되어 있어야 합니다.

## 🔍 모니터링 동작 확인

### Prometheus에서 메트릭 수집 확인

1. **포트 포워딩 시작**:
```bash
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
```

2. **브라우저에서 접근**:
   - URL: http://localhost:9090
   - Status > Targets 메뉴로 이동
   - `node-exporter` 타겟이 `UP` 상태인지 확인

3. **메트릭 쿼리 테스트**:
   - Prometheus UI에서 다음 쿼리 실행:
   ```
   up{job="node-exporter"}
   ```
   - 결과가 `1`이면 정상적으로 메트릭을 수집하고 있습니다.

### Grafana 대시보드 확인

1. **포트 포워딩 시작**:
```bash
kubectl port-forward -n monitoring svc/grafana 3000:80
```

2. **브라우저에서 접근**:
   - URL: http://localhost:3000
   - 로그인 정보:
     - Username: `admin`
     - Password: `prom-operator`

3. **대시보드 확인**:
   - 좌측 메뉴에서 "Dashboards" > "Browse" 선택
   - `node-exporter` 대시보드가 있는지 확인
   - 대시보드를 열어 노드 메트릭이 표시되는지 확인

## 🧪 자동 노드 감지 테스트 (kind 클러스터인 경우)

새 노드가 추가될 때 자동으로 node-exporter가 배포되는지 테스트:

```bash
# 새 노드 추가
kind create node --name test-node --cluster test-monitoring

# 약 30초 후 DaemonSet이 새 노드에 Pod를 생성했는지 확인
kubectl get pods -n monitoring -l app=node-exporter -o wide

# ServiceMonitor가 새 node-exporter를 자동으로 감지하는지 확인
# (Prometheus UI에서 Status > Targets 확인)
```

## 📊 메트릭 확인 명령어

### Prometheus API를 통한 확인
```bash
# 타겟 상태 확인
kubectl exec -n monitoring prometheus-prometheus-kube-prometheus-prometheus-0 -c prometheus -- \
  wget -qO- http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | select(.job=="node-exporter") | {job, health, lastScrape}'

# 메트릭 쿼리
kubectl exec -n monitoring prometheus-prometheus-kube-prometheus-prometheus-0 -c prometheus -- \
  wget -qO- "http://localhost:9090/api/v1/query?query=up{job=\"node-exporter\"}"
```

### node-exporter 메트릭 직접 확인
```bash
# node-exporter Pod에서 메트릭 확인
NODE_EXPORTER_POD=$(kubectl get pods -n monitoring -l app=node-exporter -o jsonpath='{.items[0].metadata.name}')
kubectl exec -n monitoring $NODE_EXPORTER_POD -- wget -qO- http://localhost:9100/metrics | head -20
```

## 🎯 성공 기준

다음 조건들이 모두 만족되면 설치가 성공적으로 완료된 것입니다:

- [x] 모든 Pod가 `Running` 상태
- [x] `node-exporter` DaemonSet이 모든 노드에 배포됨
- [x] `node-exporter` ServiceMonitor가 생성됨
- [ ] Prometheus에서 `node-exporter` 타겟이 `UP` 상태
- [ ] Prometheus에서 `up{job="node-exporter"}` 쿼리 결과가 `1`
- [ ] Grafana에서 node-exporter 대시보드가 표시됨
- [ ] (선택) 새 노드 추가 시 자동으로 node-exporter가 배포됨

## 🐛 문제 해결

### Prometheus에서 node-exporter를 찾을 수 없는 경우

```bash
# ServiceMonitor 확인
kubectl describe servicemonitor node-exporter -n monitoring

# Service 확인
kubectl get svc -n monitoring node-exporter -o yaml

# Endpoints 확인
kubectl get endpoints -n monitoring node-exporter
```

### Grafana 대시보드가 보이지 않는 경우

```bash
# ConfigMap 확인
kubectl get configmap -n monitoring | grep grafana

# Grafana Pod 로그 확인
kubectl logs -n monitoring -l app.kubernetes.io/name=grafana
```

