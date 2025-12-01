# 빠른 테스트 가이드

## Docker 권한 확인

다음 명령으로 Docker 접근을 확인하세요:

```bash
docker ps
```

권한 오류가 나면:

```bash
# 그룹 변경사항 적용 (현재 세션에서)
newgrp docker

# 또는 로그아웃 후 재로그인
```

## 테스트 실행

### 1. kind 클러스터 생성

```bash
kind create cluster --name test-monitoring
```

### 2. 클러스터 확인

```bash
kubectl get nodes
kubectl cluster-info
```

### 3. 설치 스크립트 실행

```bash
cd /home/yujin/promethus-example
./scripts/install.sh
```

### 4. 설치 확인

```bash
# Pod 상태 확인
kubectl get pods -n monitoring

# Service 확인
kubectl get svc -n monitoring

# DaemonSet 확인
kubectl get daemonset -n monitoring

# ServiceMonitor 확인
kubectl get servicemonitor -n monitoring
```

### 5. 접근 테스트

```bash
# Grafana 접근
kubectl port-forward -n monitoring svc/grafana 3000:80

# Prometheus 접근
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
```

