# Contracts: Kubernetes 자동 모니터링 시스템

이 디렉토리는 시스템이 생성하고 관리하는 Kubernetes 리소스의 계약(스키마)을 정의합니다.

## 계약 파일 목록

### 1. `prometheus-servicemonitor.yaml`
- **목적**: node-exporter를 위한 ServiceMonitor 리소스 스키마
- **생성 시점**: 새 노드가 감지되고 node-exporter가 배포된 후
- **관리 주체**: Controller
- **의존성**: Prometheus Operator가 설치되어 있어야 함

### 2. `node-exporter-daemonset.yaml`
- **목적**: node-exporter DaemonSet 리소스 스키마
- **생성 시점**: 첫 번째 노드가 감지되었을 때 (한 번만 생성, 모든 노드에 적용)
- **관리 주체**: Controller
- **특징**: DaemonSet이므로 새 노드 추가 시 자동으로 Pod 생성

### 3. `prometheus-rule.yaml`
- **목적**: Prometheus 알림 규칙 정의
- **생성 시점**: 모니터링 시스템 초기 설치 시
- **관리 주체**: Controller 또는 Helm Chart
- **기본 규칙**: NodeDown, HighCPUUsage, HighMemoryUsage, DiskSpaceLow

## 계약 검증

각 계약은 다음을 보장해야 합니다:

1. **스키마 유효성**: Kubernetes 리소스 스키마를 준수
2. **레이블 일관성**: 리소스 간 선택자(selector)와 레이블이 일치
3. **네임스페이스 일관성**: 모든 리소스가 동일한 네임스페이스에 배포
4. **의존성 순서**: 리소스 생성 순서가 의존성을 고려

## 테스트

계약 테스트는 다음을 검증합니다:

- 리소스 생성 시 스키마 유효성
- 레이블 및 선택자 일치
- 필수 필드 존재
- 값 범위 및 형식 검증

