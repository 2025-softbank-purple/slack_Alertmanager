# 자동 확장 테스트 가이드

## 개요

이 시스템은 **DaemonSet**을 사용하여 모든 노드에 자동으로 `node-exporter`를 배포합니다. 새 노드가 추가되면 Kubernetes가 자동으로 해당 노드에 Pod를 생성합니다.

## 자동 확장 확인 방법

### 방법 1: 자동 테스트 스크립트 사용 (권장)

```bash
./test-auto-scaling.sh
```

이 스크립트는:
1. 현재 노드와 Pod 상태를 확인
2. 새 워커 노드를 kind 클러스터에 추가
3. DaemonSet이 자동으로 새 노드에 Pod를 배포하는지 확인
4. Service Endpoints가 업데이트되는지 확인

### 방법 2: 수동 테스트

#### 1단계: 현재 상태 확인

```bash
# 현재 노드 확인
kubectl get nodes

# 현재 node-exporter Pod 확인 (노드별로 표시)
kubectl get pods -n monitoring -l app=node-exporter -o wide
```

**예상 결과**: 현재 노드 수만큼 Pod가 있어야 합니다.

#### 2단계: 새 노드 추가

```bash
# kind 클러스터에 새 워커 노드 추가
kind create node --name test-worker --cluster test-monitoring
```

#### 3단계: 노드 Ready 상태 확인

```bash
# 새 노드가 Ready 상태가 될 때까지 대기
kubectl wait --for=condition=Ready node/test-monitoring-test-worker --timeout=60s

# 노드 목록 확인
kubectl get nodes
```

**예상 결과**: 새 노드(`test-monitoring-test-worker`)가 `Ready` 상태로 표시되어야 합니다.

#### 4단계: DaemonSet 자동 배포 확인

```bash
# 약 10-30초 후 Pod 상태 확인
kubectl get pods -n monitoring -l app=node-exporter -o wide
```

**예상 결과**: 
- 새 노드에 `node-exporter-xxxxx` Pod가 자동으로 생성됨
- Pod 상태가 `Running`이어야 함
- `NODE` 컬럼에 새 노드 이름이 표시됨

#### 5단계: Service Endpoints 확인

```bash
# Service가 모든 node-exporter 인스턴스를 감지하는지 확인
kubectl get endpoints -n monitoring node-exporter
```

**예상 결과**: 
- `ENDPOINTS`에 여러 IP:포트가 표시됨 (각 노드의 node-exporter)
- 새 노드의 IP도 포함되어야 함

#### 6단계: Prometheus 자동 감지 확인

```bash
# Prometheus 포트 포워딩
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
```

브라우저에서 http://localhost:9090/targets 접속:
- `node-exporter` 타겟이 **2개 이상** 표시되어야 함 (노드 수만큼)
- 모든 타겟이 `UP` 상태여야 함

또는 API로 확인:

```bash
# Prometheus 타겟 개수 확인
kubectl exec -n monitoring prometheus-prometheus-kube-prometheus-prometheus-0 -c prometheus -- \
  wget -qO- http://localhost:9090/api/v1/targets | \
  grep -o '"job":"node-exporter"' | wc -l
```

**예상 결과**: 노드 수와 동일한 개수가 나와야 합니다.

## 실시간 모니터링

### Watch 명령으로 실시간 확인

```bash
# Pod 상태를 실시간으로 확인
kubectl get pods -n monitoring -l app=node-exporter -o wide -w

# 별도 터미널에서 새 노드 추가
# kind create node --name test-worker --cluster test-monitoring
```

**예상 동작**: 새 노드가 추가되면 자동으로 새 Pod가 생성되는 것을 실시간으로 볼 수 있습니다.

### DaemonSet 이벤트 확인

```bash
# DaemonSet 이벤트 확인
kubectl describe daemonset node-exporter -n monitoring
```

**확인 사항**:
- `Desired Number of Nodes Scheduled`: 노드 수와 일치해야 함
- `Number of Nodes Scheduled with Pods`: 모든 노드에 Pod가 배포되어야 함
- `Number Ready`: 모든 Pod가 Ready 상태여야 함

## 성공 기준

다음 조건들이 모두 만족되면 자동 확장이 정상 작동하는 것입니다:

- [ ] 새 노드 추가 후 30초 이내에 node-exporter Pod가 자동 생성됨
- [ ] 새 Pod가 `Running` 상태가 됨
- [ ] Service Endpoints에 새 노드의 IP가 추가됨
- [ ] Prometheus에서 새 타겟이 자동으로 감지됨
- [ ] 모든 타겟이 `UP` 상태로 표시됨

## 정리

테스트 완료 후 새 노드 제거:

```bash
kind delete node test-monitoring-test-worker --cluster test-monitoring
```

## 문제 해결

### 새 노드에 Pod가 생성되지 않는 경우

```bash
# DaemonSet 상태 확인
kubectl describe daemonset node-exporter -n monitoring

# 노드에 taint가 있는지 확인
kubectl describe node test-monitoring-test-worker

# Pod 이벤트 확인
kubectl get events -n monitoring --sort-by='.lastTimestamp' | grep node-exporter
```

### Service가 새 Endpoint를 감지하지 않는 경우

```bash
# Service 상세 확인
kubectl describe svc node-exporter -n monitoring

# Endpoints 상세 확인
kubectl describe endpoints node-exporter -n monitoring
```

### Prometheus가 새 타겟을 감지하지 않는 경우

```bash
# ServiceMonitor 확인
kubectl describe servicemonitor node-exporter -n monitoring

# Prometheus Operator 로그 확인
kubectl logs -n monitoring -l app.kubernetes.io/name=prometheus-operator --tail=50
```

