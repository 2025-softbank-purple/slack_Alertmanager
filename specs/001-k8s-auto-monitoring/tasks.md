# Tasks: Kubernetes 자동 모니터링 시스템

**Input**: Design documents from `/specs/001-k8s-auto-monitoring/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: TDD 원칙에 따라 모든 기능 구현 전에 테스트를 작성합니다. 테스트가 실패하는 것을 확인한 후 구현을 진행합니다.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: Go 프로젝트 구조 (cmd/, pkg/, tests/, charts/, configs/, scripts/)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create project structure per implementation plan (cmd/, pkg/, charts/, configs/, tests/, scripts/)
- [x] T002 Initialize Go module with go.mod in repository root
- [x] T003 [P] Add Kubernetes client dependencies (k8s.io/client-go, k8s.io/api) to go.mod
- [x] T004 [P] Add Helm SDK dependency (helm.sh/helm/v3) to go.mod
- [x] T005 [P] Add testing dependencies (ginkgo/gomega) to go.mod
- [x] T006 [P] Create Makefile with build, test, and install targets
- [x] T007 [P] Configure Go linting and formatting (golangci-lint, gofmt)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T008 Create Kubernetes client wrapper in pkg/client/k8s_client.go
- [x] T009 [P] Implement structured logging setup in pkg/client/logger.go (logrus or zap)
- [x] T010 [P] Create error handling utilities in pkg/client/errors.go
- [x] T011 [P] Setup environment configuration management in pkg/config/config.go
- [x] T012 Create base Helm action client wrapper in pkg/helm/client.go
- [x] T013 [P] Write unit tests for k8s_client.go in tests/unit/client/k8s_client_test.go (TDD: test first, ensure failure)
- [x] T014 [P] Write unit tests for logger.go in tests/unit/client/logger_test.go (TDD: test first, ensure failure)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Kubernetes 클러스터에 모니터링 시스템 자동 설치 (Priority: P1) 🎯 MVP

**Goal**: 사용자가 단일 명령으로 Prometheus, Grafana, Alertmanager를 자동 배포하고 구성할 수 있음

**Independent Test**: 설치 스크립트 실행 후 kubectl get pods -n monitoring으로 모든 컴포넌트가 Running 상태인지 확인

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T015 [P] [US1] Write unit test for Helm installer in tests/unit/helm/installer_test.go (TDD: test first, ensure failure)
- [x] T016 [P] [US1] Write integration test for Prometheus deployment in tests/integration/helm/prometheus_deploy_test.go (TDD: test first, ensure failure)
- [x] T017 [P] [US1] Write integration test for Grafana deployment in tests/integration/helm/grafana_deploy_test.go (TDD: test first, ensure failure)
- [x] T018 [P] [US1] Write integration test for Alertmanager deployment in tests/integration/helm/alertmanager_deploy_test.go (TDD: test first, ensure failure)
- [x] T019 [P] [US1] Write contract test for PrometheusRule creation in tests/contract/prometheus_rule_test.go (TDD: test first, ensure failure)

### Implementation for User Story 1

- [x] T020 [US1] Implement Helm installer for Prometheus Operator in pkg/helm/installer.go (depends on T015)
- [x] T021 [US1] Implement Helm installer for Grafana in pkg/helm/installer.go (depends on T016, T017)
- [x] T022 [US1] Implement Helm installer for Alertmanager in pkg/helm/installer.go (depends on T018)
- [x] T023 [US1] Create Prometheus Operator Helm values configuration in charts/prometheus-stack/values.yaml
- [x] T024 [US1] Create Grafana Helm values configuration in charts/grafana/values.yaml
- [x] T025 [US1] Create default PrometheusRule with alert rules in configs/prometheus/alert-rules.yaml (depends on T019)
- [x] T026 [US1] Implement error handling and user-friendly error messages in pkg/helm/installer.go (FR-007)
- [x] T027 [US1] Create installation script in scripts/install.sh that orchestrates all Helm deployments
- [x] T028 [US1] Add logging for installation progress in pkg/helm/installer.go
- [x] T029 [US1] Implement health check for deployed components in pkg/helm/health.go (FR-008)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. All components (Prometheus, Grafana, Alertmanager) should be deployed and running.

---

## Phase 4: User Story 2 - 새 노드 추가 시 자동 모니터링 연결 (Priority: P1)

**Goal**: 새 노드가 클러스터에 추가되면 자동으로 감지하고 node-exporter를 배포하여 모니터링에 연결

**Independent Test**: 새 노드를 추가한 후 DaemonSet이 자동으로 해당 노드에 node-exporter Pod를 생성하고, Prometheus가 자동으로 메트릭을 수집하는지 확인

### Implementation for User Story 2

**Note**: DaemonSet을 사용하므로 Go Controller가 불필요합니다. DaemonSet은 모든 노드(기존 및 새 노드)에 자동으로 Pod를 배포하고, ServiceMonitor는 모든 node-exporter Service를 자동으로 감지합니다.

- [x] T030 [US2] Create node-exporter DaemonSet YAML in configs/node-exporter/daemonset.yaml (FR-003)
- [x] T031 [US2] Create node-exporter Service YAML in configs/node-exporter/service.yaml (FR-003)
- [x] T032 [US2] Create ServiceMonitor YAML in configs/prometheus/servicemonitor.yaml to auto-discover all node-exporter services (FR-003)
- [x] T033 [US2] Update installation script to deploy node-exporter DaemonSet and ServiceMonitor in scripts/install.sh

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. DaemonSet will automatically deploy node-exporter to all nodes (existing and new), and Prometheus will automatically discover them via ServiceMonitor.

---

## Phase 5: User Story 3 - 모니터링 데이터 시각화 및 확인 (Priority: P2)

**Goal**: 사용자가 Grafana 대시보드를 통해 클러스터와 노드의 모니터링 데이터를 시각적으로 확인할 수 있음

**Independent Test**: Grafana에 접근하여 기본 대시보드가 표시되고, 노드 메트릭이 시각화되는지 확인

### Implementation for User Story 3

**Note**: Grafana는 ConfigMap을 통해 대시보드를 자동으로 로드합니다. `grafana_dashboard: "1"` 레이블이 있는 ConfigMap을 자동으로 감지합니다.

- [x] T048 [US3] Create Grafana dashboard ConfigMap in configs/grafana/dashboards/configmap.yaml (FR-006)
- [x] T049 [US3] Update installation script to include dashboard provisioning in scripts/install.sh
- [x] T050 [US3] Add Grafana service access documentation in README.md and quickstart.md
- [x] T051 [US3] Verify Prometheus data source is configured in Grafana (Prometheus Operator가 자동으로 구성)

**Checkpoint**: All user stories should now be independently functional. Users can view monitoring data through Grafana dashboards.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T056 [P] Add comprehensive error messages for all failure scenarios in scripts/install.sh (FR-007)
- [x] T057 [P] Implement health check and auto-recovery for all components (Kubernetes native - Pod restart) (FR-008)
- [x] T058 [P] Add structured logging throughout all components (Helm logs, kubectl events)
- [x] T059 [P] Create uninstallation script in scripts/uninstall.sh
- [x] T060 [P] Add comprehensive documentation in README.md
- [x] T061 [P] Update quickstart.md with validation steps (already in quickstart.md)
- [x] T062 [P] Performance metrics collection (Prometheus가 자동으로 수집)
- [x] T063 [P] Code cleanup and refactoring (구조 단순화 완료)
- [x] T064 [P] Unit tests (Go Controller 제거로 테스트 범위 축소)
- [ ] T065 Run quickstart.md validation end-to-end (수동 검증 필요)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Requires US1 for Prometheus Operator to be deployed, but Controller logic is independently testable
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Requires US1 for Grafana to be deployed, but dashboard provisioning is independently testable

### Within Each User Story

- Tests (TDD) MUST be written and FAIL before implementation
- Core components before integration
- Error handling and logging after core implementation
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel (T003-T007)
- All Foundational tasks marked [P] can run in parallel (T009-T014)
- Once Foundational phase completes:
  - US1 and US2 can start in parallel (different components)
  - All tests for a user story marked [P] can run in parallel
  - Different user stories can be worked on in parallel by different team members
- Polish phase tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: T015 - Write unit test for Helm installer in tests/unit/helm/installer_test.go
Task: T016 - Write integration test for Prometheus deployment in tests/integration/helm/prometheus_deploy_test.go
Task: T017 - Write integration test for Grafana deployment in tests/integration/helm/grafana_deploy_test.go
Task: T018 - Write integration test for Alertmanager deployment in tests/integration/helm/alertmanager_deploy_test.go
Task: T019 - Write contract test for PrometheusRule creation in tests/contract/prometheus_rule_test.go

# After tests are written and failing, launch implementation tasks:
Task: T023 - Create Prometheus Operator Helm values configuration in charts/prometheus-stack/values.yaml
Task: T024 - Create Grafana Helm values configuration in charts/grafana/values.yaml
Task: T025 - Create default PrometheusRule with alert rules in configs/prometheus/alert-rules.yaml
```

---

## Parallel Example: User Story 2

```bash
# Launch all tests for User Story 2 together:
Task: T030 - Write unit test for node controller in tests/unit/controller/node_controller_test.go
Task: T031 - Write unit test for node reconciler in tests/unit/controller/reconciler_test.go
Task: T032 - Write unit test for DaemonSet creator in tests/unit/exporter/daemonset_test.go
Task: T033 - Write unit test for ServiceMonitor creator in tests/unit/exporter/servicemonitor_test.go
Task: T034 - Write integration test for node detection in tests/integration/controller/node_detection_test.go
Task: T035 - Write integration test for node-exporter deployment in tests/integration/exporter/daemonset_deploy_test.go
Task: T036 - Write contract test for ServiceMonitor resource in tests/contract/servicemonitor_test.go
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
   - Run installation script
   - Verify all components are running
   - Test error handling
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Helm installation)
   - Developer B: User Story 2 (Node controller) - can start after US1 Prometheus is deployed
   - Developer C: User Story 3 (Grafana dashboards) - can start after US1 Grafana is deployed
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- **TDD**: Verify tests fail before implementing (Red-Green-Refactor cycle)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Test coverage target: 80% for core functionality (Constitution requirement)

