# Research: Kubernetes 자동 모니터링 시스템

**Date**: 2024-12-01  
**Feature**: 001-k8s-auto-monitoring

## Research Topics

### 1. Kubernetes Controller 구현 패턴

**Decision**: controller-runtime 프레임워크를 사용한 Reconcile 패턴 구현

**Rationale**:
- Kubernetes의 표준 Controller 패턴을 따름
- Reconcile 루프를 통해 원하는 상태와 실제 상태를 지속적으로 동기화
- 이벤트 기반 처리로 리소스 효율적
- 테스트 가능한 구조 (fake client 사용 가능)

**Alternatives considered**:
- 직접 informer 사용: 더 많은 보일러플레이트 코드 필요
- Operator SDK: 과도한 기능, 단순 Controller에는 불필요
- 폴링 방식: 리소스 비효율적, 실시간 감지 불가

**Implementation approach**:
- `sigs.k8s.io/controller-runtime` 패키지 사용
- Node 리소스를 Watch하여 추가/삭제 이벤트 감지
- Reconcile 함수에서 노드 상태 확인 및 필요한 액션 수행

---

### 2. Helm Chart 배포 및 관리

**Decision**: Helm Go SDK를 사용하여 프로그래밍 방식으로 Chart 배포

**Rationale**:
- Helm CLI보다 더 세밀한 제어 가능
- 에러 처리 및 상태 확인 용이
- Go 프로젝트와 자연스럽게 통합
- Chart 업그레이드 및 롤백 지원

**Alternatives considered**:
- Helm CLI 실행: 외부 프로세스 의존, 에러 처리 복잡
- 직접 YAML 적용: Chart의 의존성 관리 어려움
- Kustomize: Helm Chart와의 통합 복잡

**Implementation approach**:
- `helm.sh/helm/v3/pkg/action` 패키지 사용
- Install/Upgrade 액션으로 Chart 배포
- Values 파일을 동적으로 생성하여 설정 커스터마이징

---

### 3. Prometheus Operator 및 ServiceMonitor

**Decision**: Prometheus Operator Helm Chart 사용, ServiceMonitor로 자동 스크랩 설정

**Rationale**:
- Prometheus Operator는 Prometheus 생명주기 관리 자동화
- ServiceMonitor CRD를 통해 자동으로 스크랩 대상 추가/제거
- 노드 추가 시 ServiceMonitor만 생성하면 Prometheus가 자동으로 감지
- 표준 Kubernetes 모니터링 패턴

**Alternatives considered**:
- 직접 Prometheus 설정: 수동 reload 필요, 복잡한 설정 관리
- Prometheus 파일 기반 서비스 디스커버리: 덜 자동화됨
- Static config: 노드 추가 시 수동 설정 필요

**Implementation approach**:
- `prometheus-community/kube-prometheus-stack` Helm Chart 사용
- 새 노드 감지 시 해당 노드의 node-exporter를 위한 ServiceMonitor 생성
- Prometheus Operator가 자동으로 ServiceMonitor를 감지하여 스크랩 시작

---

### 4. node-exporter DaemonSet 배포

**Decision**: DaemonSet으로 각 노드에 node-exporter 자동 배포

**Rationale**:
- DaemonSet은 클러스터의 모든 노드에 자동으로 Pod 배포
- 새 노드 추가 시 자동으로 node-exporter Pod 생성
- 노드 제거 시 Pod도 자동으로 정리
- Kubernetes 네이티브 방식

**Alternatives considered**:
- 수동 배포: 자동화 요구사항 미충족
- StatefulSet: 노드별 고유성이 불필요
- Deployment: 모든 노드에 배포 보장 불가

**Implementation approach**:
- Controller가 새 노드 감지 시 해당 노드에 node-exporter DaemonSet 배포
- Service 리소스로 node-exporter 엔드포인트 노출
- ServiceMonitor로 Prometheus 스크랩 설정

---

### 5. Grafana 대시보드 자동 설치

**Decision**: Grafana Dashboard ConfigMap을 통한 자동 설치

**Rationale**:
- Grafana는 ConfigMap을 통해 대시보드를 자동으로 로드
- `grafana_dashboard` 레이블이 있는 ConfigMap을 자동으로 감지
- 대시보드 JSON을 ConfigMap에 포함하여 배포 시 자동 설치
- 표준 Grafana 운영 패턴

**Alternatives considered**:
- Grafana API 사용: 인증 및 네트워크 의존성
- 수동 import: 자동화 요구사항 미충족
- Provisioning 파일: ConfigMap보다 복잡

**Implementation approach**:
- Kubernetes/Prometheus 공식 대시보드 JSON 다운로드
- ConfigMap으로 생성하여 Grafana 네임스페이스에 배포
- `grafana_dashboard: "1"` 레이블 추가하여 자동 로드

---

### 6. Alertmanager 기본 알림 규칙

**Decision**: PrometheusRule CRD를 통한 알림 규칙 정의

**Rationale**:
- Prometheus Operator는 PrometheusRule CRD를 자동으로 감지
- Kubernetes 네이티브 방식으로 알림 규칙 관리
- GitOps와 호환
- 규칙 업데이트 시 자동으로 Prometheus에 반영

**Alternatives considered**:
- Prometheus 설정 파일: 수동 reload 필요
- Alertmanager 설정: Prometheus와 분리되어 관리 복잡
- Grafana 알림: Prometheus 알림 규칙과 중복

**Implementation approach**:
- 기본 알림 규칙 정의:
  - NodeDown: 노드가 다운된 경우
  - HighCPUUsage: CPU 사용률이 80% 이상인 경우
  - HighMemoryUsage: 메모리 사용률이 80% 이상인 경우
  - DiskSpaceLow: 디스크 공간이 20% 미만인 경우
- PrometheusRule 리소스로 생성
- Prometheus Operator가 자동으로 로드

---

### 7. 에러 처리 및 복구 전략

**Decision**: Kubernetes의 자동 복구 메커니즘 활용 + Controller 레벨 재시도

**Rationale**:
- Kubernetes는 Pod 실패 시 자동으로 재시작
- Deployment/DaemonSet의 replicas 보장
- Controller는 Reconcile 실패 시 exponential backoff로 재시도
- 최종 일관성 보장

**Implementation approach**:
- Controller Reconcile에서 에러 발생 시 에러 반환하여 재시도 스케줄링
- Helm 배포 실패 시 명확한 에러 메시지 로깅
- 노드 배포 실패 시 이벤트 생성하여 사용자에게 알림
- Health check를 통한 구성 요소 상태 모니터링

---

### 8. 테스트 전략

**Decision**: 
- 단위 테스트: fake Kubernetes client 사용
- 통합 테스트: kind (Kubernetes in Docker) 클러스터 사용

**Rationale**:
- fake client로 Controller 로직 독립 테스트 가능
- kind로 실제 Kubernetes 환경과 유사한 통합 테스트
- CI/CD 파이프라인에서 실행 가능
- 빠른 피드백 루프

**Implementation approach**:
- `sigs.k8s.io/controller-runtime/pkg/client/fake` 사용
- ginkgo/gomega로 BDD 스타일 테스트 작성
- testcontainers 또는 kind로 통합 테스트 환경 구성
- 각 컴포넌트별 독립적인 테스트 스위트

---

## 기술 스택 최종 결정

| 컴포넌트 | 기술 선택 | 버전 |
|---------|---------|------|
| 언어 | Go | 1.21+ |
| Kubernetes 클라이언트 | controller-runtime | v0.16.0+ |
| Helm SDK | helm.sh/helm/v3 | v3.12.0+ |
| Prometheus | prometheus-community/kube-prometheus-stack | v55.0.0+ |
| Grafana | grafana/grafana | v6.50.0+ |
| 테스트 프레임워크 | ginkgo/gomega | v2.0+ |
| 통합 테스트 | kind | v0.20.0+ |

## 참고 자료

- [Kubernetes Controller Pattern](https://kubernetes.io/docs/concepts/architecture/controller/)
- [controller-runtime Documentation](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
- [Helm Go SDK](https://pkg.go.dev/helm.sh/helm/v3)
- [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator)
- [Grafana Dashboard Provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/)

