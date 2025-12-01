# Implementation Plan: Kubernetes 자동 모니터링 시스템

**Branch**: `001-k8s-auto-monitoring` | **Date**: 2024-12-01 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-k8s-auto-monitoring/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Kubernetes 클러스터에 Prometheus, Grafana, Alertmanager를 자동으로 배포하고, 새 노드가 추가될 때 자동으로 감지하여 node-exporter를 배포하고 모니터링에 연결하는 시스템을 구축합니다. Helm Chart를 통한 배포와 Kubernetes API Watch를 통한 실시간 노드 감지가 핵심입니다.

## Technical Context

**Language/Version**: Go 1.21+ (Kubernetes 클라이언트 라이브러리 및 Controller 구현)  
**Primary Dependencies**: 
- k8s.io/client-go (Kubernetes API 클라이언트)
- k8s.io/api (Kubernetes 리소스 정의)
- helm.sh/helm/v3 (Helm Chart 배포)
- prometheus-operator (Prometheus Operator Helm Chart)
- grafana/grafana (Grafana Helm Chart)

**Storage**: 
- Prometheus: PVC (PersistentVolumeClaim)를 통한 시계열 데이터 저장
- Grafana: ConfigMap을 통한 대시보드 및 설정 저장

**Testing**: 
- Go testing 패키지 (단위 테스트)
- testcontainers 또는 kind (통합 테스트용 Kubernetes 클러스터)
- ginkgo/gomega (BDD 스타일 테스트)

**Target Platform**: Linux (로컬 Kubernetes 환경: minikube, kind, k3s 등)  
**Project Type**: single (단일 프로젝트, Kubernetes Controller/Operator)  
**Performance Goals**: 
- 노드 감지: 2분 이내 (SC-002)
- 설치 시간: 5분 이내 (SC-001)
- 메트릭 수집 가용성: 95% 이상 (SC-003)

**Constraints**: 
- 로컬 Kubernetes 환경에서만 작동
- 최소 리소스: 2GB RAM, 2 CPU 코어
- 네트워크 정책이 모니터링 트래픽 허용 필요

**Scale/Scope**: 
- 최대 10개 노드까지 지원 (로컬 환경 가정)
- 단일 클러스터 모니터링

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**TDD Compliance**:
- [x] 모든 기능 구현 전에 테스트 작성 계획이 명확한가? → ✅ Controller 로직, Helm 배포, 노드 감지 각각에 대해 테스트 작성 계획 수립됨 (research.md 참조)
- [x] 테스트가 독립적으로 실행 가능하도록 설계되었는가? → ✅ 각 컴포넌트별 모킹 및 테스트 더블 사용, fake Kubernetes client 활용
- [x] 통합 테스트가 필요한 영역(서비스 간 통신, Kubernetes 리소스, Prometheus/Grafana 연동)이 식별되었는가? → ✅ Kubernetes 리소스 생성/관리, Prometheus ServiceMonitor, Grafana 대시보드 연동 테스트 계획 수립됨

**Simplicity Check**:
- [x] 현재 요구사항에 집중하고 있는가? (YAGNI 원칙) → ✅ 로컬 환경만 지원, 기본 기능에 집중, 미래 확장성 고려하지 않음
- [x] 불필요한 복잡성이 없는가? → ✅ 기존 Helm Chart 활용, 표준 Kubernetes 패턴 사용, 커스텀 Operator는 최소화
- [x] 모든 복잡성은 정당화되었는가? → ✅ Kubernetes API Watch는 자동 감지 요구사항에 필수 (Complexity Tracking 참조)

**Observability Check**:
- [x] 로깅 전략이 정의되었는가? → ✅ 구조화된 로깅 (logrus 또는 zap), 노드 감지, 배포 상태, 에러 로깅 계획 수립됨
- [x] 에러 처리 및 메트릭 수집 계획이 있는가? → ✅ Controller 메트릭 수집, Prometheus로 노출, 에러 이벤트 생성 계획 수립됨

**Post-Design Status**: ✅ 모든 체크리스트 항목 통과. 설계 단계에서 TDD, 단순성, 관찰 가능성 원칙을 준수함을 확인.

## Project Structure

### Documentation (this feature)

```text
specs/001-k8s-auto-monitoring/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
.
├── cmd/
│   └── node-watcher/          # 노드 감지 Controller 메인
│       └── main.go
├── pkg/
│   ├── controller/            # Kubernetes Controller 로직
│   │   ├── node_controller.go
│   │   └── reconciler.go
│   ├── helm/                  # Helm Chart 배포 로직
│   │   ├── installer.go
│   │   └── config.go
│   ├── exporter/              # node-exporter 배포 로직
│   │   ├── daemonset.go
│   │   └── servicemonitor.go
│   └── client/                # Kubernetes 클라이언트 래퍼
│       └── k8s_client.go
├── charts/                    # Helm Chart 정의
│   ├── prometheus-stack/      # Prometheus Operator Chart
│   ├── grafana/               # Grafana Chart
│   └── node-watcher/          # 노드 감지 Controller Chart
├── configs/                   # 설정 파일
│   ├── prometheus/
│   │   ├── alert-rules.yaml
│   │   └── servicemonitor.yaml
│   └── grafana/
│       └── dashboards/
├── tests/
│   ├── unit/                  # 단위 테스트
│   │   ├── controller/
│   │   ├── helm/
│   │   └── exporter/
│   ├── integration/           # 통합 테스트
│   │   ├── k8s_resources_test.go
│   │   └── helm_deploy_test.go
│   └── contract/              # 계약 테스트
│       └── api_test.go
├── scripts/                   # 설치/배포 스크립트
│   ├── install.sh
│   └── uninstall.sh
├── go.mod
├── go.sum
└── Makefile
```

**Structure Decision**: 단일 Go 프로젝트로 구성. Kubernetes Controller 패턴을 사용하여 노드 감지 및 자동 배포를 구현. Helm Chart는 charts/ 디렉토리에 배치하여 재사용 가능하도록 구성.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Kubernetes Controller 구현 | 실시간 노드 감지 및 자동 배포 요구사항 | 단순 스크립트는 실시간 감지 불가, 폴링은 리소스 비효율적 |
| Helm Chart 통합 | Prometheus/Grafana 표준 배포 방식 | 직접 YAML 배포는 업그레이드/관리 복잡도 증가 |
