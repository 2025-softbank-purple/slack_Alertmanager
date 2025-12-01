# 테스트 가이드

## 전제 조건 확인

Kubernetes 클러스터가 실행 중이어야 합니다. 다음 중 하나를 사용할 수 있습니다:

### 옵션 1: kind 사용 (권장 - 로컬 테스트용)

```bash
# kind 설치 (sudo 권한 필요)
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# 테스트 클러스터 생성
kind create cluster --name test-monitoring

# kubectl이 kind 클러스터를 사용하도록 설정
kubectl cluster-info --context kind-test-monitoring
```

### 옵션 2: minikube 사용

```bash
# minikube 설치
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# 클러스터 시작
minikube start

# kubectl이 minikube를 사용하도록 설정
kubectl config use-context minikube
```

### 옵션 3: 기존 클러스터 사용

이미 Kubernetes 클러스터가 있다면:

```bash
# 클러스터 연결 확인
kubectl cluster-info
kubectl get nodes
```

## 설치 테스트

### 1. 클러스터 상태 확인

```bash
kubectl get nodes
kubectl cluster-info
```

### 2. 설치 스크립트 실행

```bash
cd /home/yujin/promethus-example
./scripts/install.sh
```

### 3. 설치 확인

```bash
# 네임스페이스 확인
kubectl get namespace monitoring

# Pod 상태 확인
kubectl get pods -n monitoring

# 예상 출력:
# NAME                                    READY   STATUS    RESTARTS   AGE
# prometheus-prometheus-0                 2/2     Running   0          2m
# grafana-xxx                             1/1     Running   0          2m
# alertmanager-xxx                         1/1     Running   0          2m
# node-exporter-xxx                        1/1     Running   0          1m

# Service 확인
kubectl get svc -n monitoring

# ServiceMonitor 확인
kubectl get servicemonitor -n monitoring

# PrometheusRule 확인
kubectl get prometheusrule -n monitoring
```

### 4. node-exporter DaemonSet 확인

```bash
# DaemonSet 확인
kubectl get daemonset -n monitoring

# 모든 노드에 Pod가 배포되었는지 확인
kubectl get pods -n monitoring -l app=node-exporter -o wide
```

### 5. 새 노드 추가 테스트 (kind의 경우)

```bash
# kind에 새 노드 추가
kind create node --name test-node --cluster test-monitoring

# DaemonSet이 새 노드에 자동으로 Pod를 생성했는지 확인 (약 30초 후)
kubectl get pods -n monitoring -l app=node-exporter -o wide

# ServiceMonitor가 모든 node-exporter를 감지하는지 확인
kubectl get servicemonitor node-exporter -n monitoring -o yaml
```

### 6. Prometheus 타겟 확인

```bash
# Prometheus 포트 포워딩
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090

# 브라우저에서 http://localhost:9090 접근
# Status > Targets에서 node-exporter가 스크랩되고 있는지 확인
```

### 7. Grafana 접근

```bash
# Grafana 포트 포워딩
kubectl port-forward -n monitoring svc/grafana 3000:80

# 브라우저에서 http://localhost:3000 접근
# 로그인: admin / prom-operator
# 대시보드에서 노드 메트릭 확인
```

## 문제 해결

### Pod가 시작되지 않는 경우

```bash
# Pod 이벤트 확인
kubectl describe pod <pod-name> -n monitoring

# 로그 확인
kubectl logs <pod-name> -n monitoring
```

### ServiceMonitor가 작동하지 않는 경우

```bash
# ServiceMonitor 상세 확인
kubectl describe servicemonitor node-exporter -n monitoring

# Prometheus Operator 로그 확인
kubectl logs -n monitoring -l app.kubernetes.io/name=prometheus-operator
```

### 리소스 부족 오류

```bash
# 노드 리소스 확인
kubectl top nodes

# Pod 리소스 확인
kubectl top pods -n monitoring
```

## 정리

테스트 완료 후 클러스터 정리:

```bash
# 설치 제거
./scripts/uninstall.sh

# kind 클러스터 삭제 (kind 사용 시)
kind delete cluster --name test-monitoring
```

