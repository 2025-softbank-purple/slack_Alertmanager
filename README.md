# Kubernetes 자동 모니터링 시스템

로컬 Kubernetes 클러스터에 Prometheus, Grafana, Alertmanager를 자동으로 배포하고, 새 노드가 추가될 때 자동으로 모니터링에 연결하는 시스템입니다.

## 주요 특징

- **자동 설치**: 단일 명령으로 Prometheus, Grafana, Alertmanager 배포
- **자동 노드 감지**: DaemonSet을 사용하여 모든 노드(기존 및 새 노드)에 자동으로 node-exporter 배포
- **자동 모니터링 연결**: ServiceMonitor를 통해 Prometheus가 모든 node-exporter를 자동으로 감지 및 스크랩
- **기본 대시보드**: Kubernetes/Prometheus 표준 대시보드 자동 설치
- **기본 알림 규칙**: 노드 다운, 높은 CPU/메모리 사용률 등 기본 알림 규칙 제공

## 아키텍처

이 시스템은 **Go Controller 없이** Kubernetes의 네이티브 기능만 사용합니다:

1. **DaemonSet**: node-exporter를 모든 노드에 자동 배포
   - 새 노드 추가 시 자동으로 Pod 생성
   - 노드 제거 시 Pod 자동 정리

2. **ServiceMonitor**: Prometheus가 모든 node-exporter Service를 자동 감지
   - 하나의 ServiceMonitor로 모든 node-exporter 자동 스크랩
   - 노드별 ServiceMonitor 생성 불필요

## 전제 조건

- 로컬 Kubernetes 클러스터 실행 중 (minikube, kind, k3s 등)
- `kubectl` 설치 및 클러스터 접근 가능
- `helm` v3.12+ 설치
- 최소 2GB RAM, 2 CPU 코어 사용 가능

## 빠른 시작

### 설치

```bash
# 저장소 클론
git clone <repository-url>
cd promethus-example

# 설치 스크립트 실행
./scripts/install.sh
```

### 접근

```bash
# Grafana 접근
kubectl port-forward -n monitoring svc/grafana 3000:80
# 브라우저에서 http://localhost:3000 접근
# 기본 사용자: admin, 비밀번호: prom-operator

# Prometheus 접근
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
# 브라우저에서 http://localhost:9090 접근
```

## 프로젝트 구조

```
.
├── charts/                    # Helm Chart values
│   ├── prometheus-stack/
│   └── grafana/
├── configs/                   # Kubernetes 리소스 YAML
│   ├── prometheus/
│   │   ├── alert-rules.yaml
│   │   └── servicemonitor.yaml
│   ├── node-exporter/
│   │   ├── daemonset.yaml
│   │   └── service.yaml
│   └── grafana/
│       └── dashboards/
├── scripts/                   # 설치/배포 스크립트
│   ├── install.sh
│   └── uninstall.sh
└── specs/                     # 기능 스펙 및 문서
    └── 001-k8s-auto-monitoring/
```

## 작동 원리

### 자동 노드 감지 및 모니터링

1. **DaemonSet 배포**: `node-exporter` DaemonSet을 배포하면 모든 노드에 자동으로 Pod가 생성됩니다.
2. **Service 생성**: node-exporter Service가 생성되어 모든 Pod의 엔드포인트를 노출합니다.
3. **ServiceMonitor 생성**: 하나의 ServiceMonitor가 모든 node-exporter Service를 자동으로 감지합니다.
4. **Prometheus 스크랩**: Prometheus Operator가 ServiceMonitor를 감지하고 자동으로 메트릭을 수집합니다.

**새 노드 추가 시**:
- DaemonSet이 자동으로 새 노드에 node-exporter Pod 생성
- Service가 자동으로 새 Pod를 엔드포인트에 추가
- Prometheus가 자동으로 새 엔드포인트를 스크랩 시작
- **추가 작업 불필요!**

## 제거

```bash
./scripts/uninstall.sh
```

또는 수동으로:

```bash
helm uninstall prometheus -n monitoring
helm uninstall grafana -n monitoring
kubectl delete -f configs/node-exporter/
kubectl delete -f configs/prometheus/servicemonitor.yaml
kubectl delete namespace monitoring
```

## 참고 자료

- [Prometheus 문서](https://prometheus.io/docs/)
- [Grafana 문서](https://grafana.com/docs/)
- [Kubernetes DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)
- [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator)

