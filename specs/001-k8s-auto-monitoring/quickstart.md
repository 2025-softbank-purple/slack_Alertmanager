# Quickstart: Kubernetes 자동 모니터링 시스템

**Date**: 2024-12-01  
**Feature**: 001-k8s-auto-monitoring

## 전제 조건

- 로컬 Kubernetes 클러스터 실행 중 (minikube, kind, k3s 등)
- `kubectl` 설치 및 클러스터 접근 가능
- `helm` v3.12+ 설치
- Go 1.21+ (개발 환경)
- 최소 2GB RAM, 2 CPU 코어 사용 가능

## 빠른 시작 (5분 이내)

### 1. 저장소 클론 및 빌드

```bash
# 저장소 클론
git clone <repository-url>
cd promethus-example

# 의존성 설치
go mod download

# Controller 빌드
make build
```

### 2. Kubernetes 클러스터 준비

```bash
# 클러스터 상태 확인
kubectl cluster-info
kubectl get nodes

# 모니터링 네임스페이스 생성
kubectl create namespace monitoring
```

### 3. 모니터링 시스템 설치

```bash
# 설치 스크립트 실행
./scripts/install.sh

# 또는 수동 설치
# 1. Prometheus Operator 설치
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace

# 2. Grafana 대시보드 설치
kubectl apply -f configs/grafana/dashboards/

# 3. Controller 배포
kubectl apply -f charts/node-watcher/
```

### 4. 설치 확인

```bash
# 모든 Pod가 Running 상태인지 확인
kubectl get pods -n monitoring

# 예상 출력:
# NAME                                    READY   STATUS    RESTARTS   AGE
# prometheus-prometheus-0                 2/2     Running   0          2m
# grafana-xxx                             1/1     Running   0          2m
# alertmanager-xxx                        1/1     Running   0          2m
# node-exporter-xxx                        1/1     Running   0          1m
# node-watcher-controller-xxx             1/1     Running   0          1m

# Service 확인
kubectl get svc -n monitoring

# Grafana 접근 (포트 포워딩)
kubectl port-forward -n monitoring svc/grafana 3000:80
# 브라우저에서 http://localhost:3000 접근
# 기본 사용자: admin, 비밀번호: prom-operator (또는 설치 시 설정한 값)
```

### 5. 노드 추가 테스트

```bash
# 새 노드 추가 (kind 예시)
kind create node --name test-node

# Controller 로그 확인 (노드 감지 확인)
kubectl logs -n monitoring -l app=node-watcher-controller -f

# node-exporter가 새 노드에 배포되었는지 확인
kubectl get pods -n monitoring -l app=node-exporter -o wide

# ServiceMonitor 생성 확인
kubectl get servicemonitor -n monitoring

# Prometheus 타겟 확인
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
# 브라우저에서 http://localhost:9090 접근
# Status > Targets에서 node-exporter가 스크랩되고 있는지 확인
```

## 기본 사용법

### Grafana 대시보드 확인

1. Grafana에 접근: `http://localhost:3000`
2. 로그인: admin / prom-operator (또는 설정한 비밀번호)
3. 대시보드 메뉴에서 다음 대시보드를 확인:
   - **Kubernetes / Compute Resources / Node**: 노드별 CPU, 메모리 사용률
   - **Node Exporter Full**: 노드의 상세 메트릭

### 알림 규칙 확인

```bash
# PrometheusRule 확인
kubectl get prometheusrule -n monitoring

# 알림 규칙 상세 확인
kubectl describe prometheusrule k8s-node-alerts -n monitoring

# Alertmanager UI 접근
kubectl port-forward -n monitoring svc/alertmanager 9093:9093
# 브라우저에서 http://localhost:9093 접근
```

### 노드 모니터링 상태 확인

```bash
# 현재 모니터링 중인 노드 목록
kubectl get nodes

# 각 노드의 node-exporter 상태
kubectl get pods -n monitoring -l app=node-exporter -o wide

# ServiceMonitor 목록 (노드별)
kubectl get servicemonitor -n monitoring -l app=node-exporter
```

## 문제 해결

### Pod가 시작되지 않는 경우

```bash
# Pod 이벤트 확인
kubectl describe pod <pod-name> -n monitoring

# 로그 확인
kubectl logs <pod-name> -n monitoring

# 리소스 부족 확인
kubectl top nodes
kubectl top pods -n monitoring
```

### 노드가 감지되지 않는 경우

```bash
# Controller 로그 확인
kubectl logs -n monitoring -l app=node-watcher-controller

# Controller 이벤트 확인
kubectl get events -n monitoring --field-selector involvedObject.name=node-watcher-controller

# RBAC 권한 확인
kubectl auth can-i get nodes --as=system:serviceaccount:monitoring:node-watcher
```

### Prometheus가 메트릭을 수집하지 않는 경우

```bash
# ServiceMonitor 확인
kubectl get servicemonitor -n monitoring
kubectl describe servicemonitor <name> -n monitoring

# Prometheus 타겟 확인
# http://localhost:9090/targets 에서 상태 확인

# node-exporter 엔드포인트 확인
kubectl port-forward -n monitoring svc/node-exporter 9100:9100
curl http://localhost:9100/metrics
```

## 제거

```bash
# 전체 시스템 제거
./scripts/uninstall.sh

# 또는 수동 제거
helm uninstall prometheus -n monitoring
kubectl delete namespace monitoring
```

## 다음 단계

- [ ] 커스텀 알림 규칙 추가
- [ ] Grafana 대시보드 커스터마이징
- [ ] 장기 데이터 보관 설정 (Thanos 등)
- [ ] 외부 알림 채널 설정 (Slack, Email 등)

## 참고 자료

- [Prometheus 문서](https://prometheus.io/docs/)
- [Grafana 문서](https://grafana.com/docs/)
- [Kubernetes 모니터링 가이드](https://kubernetes.io/docs/tasks/debug/debug-cluster/resource-metrics-pipeline/)

